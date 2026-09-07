package myrasec

import (
	"net/http"
	"time"
)

// responseCache holds a cached HTTP response with expiration metadata.
type responseCache struct {
	Key     string
	Created int64
	Expire  int64
	Body    any
}

// isExpired checks if the cached response is expired
func (c *responseCache) isExpired() bool {
	return c.Expire < time.Now().Unix()
}

// inCache reports whether a valid (unexpired, non-empty) response for the passed request
// is stored in the cache.
func (api *API) inCache(req *http.Request) bool {
	return api.fromCache(req) != nil
}

// fromCache returns the cached response for the passed request, or nil when there is
// none or the entry expired. An expired entry is removed. Lookup, expiry check and
// removal happen under one lock so a concurrent flush cannot slip in between.
func (api *API) fromCache(req *http.Request) any {
	s := BuildCacheKey(req)

	api.muCache.Lock()
	defer api.muCache.Unlock()

	c, ok := api.cache[s]
	if !ok {
		return nil
	}

	if c.isExpired() {
		delete(api.cache, s)
		return nil
	}

	return c.Body
}

// cacheGeneration returns the current generation of the cache. Every flush and every
// removal of an entry starts a new one.
func (api *API) cacheGeneration() uint64 {
	api.muCache.Lock()
	defer api.muCache.Unlock()

	return api.cacheGen
}

// cacheResponse stores the response body in the cache
func (api *API) cacheResponse(req *http.Request, resp any) {
	api.cacheResponseOfGeneration(req, resp, api.cacheGeneration())
}

// cacheResponseOfGeneration stores the response body in the cache, unless the cache was
// flushed since the passed generation was read. A caller reads the generation before it
// sends the request. A response that was requested before a flush can hold the state from
// before the write that caused the flush, it must not come back into the cache.
func (api *API) cacheResponseOfGeneration(req *http.Request, resp any, generation uint64) {
	if !api.caching {
		return
	}

	s := BuildCacheKey(req)
	api.muCache.Lock()
	defer api.muCache.Unlock()

	if api.cacheGen != generation {
		return
	}

	api.cache[s] = &responseCache{
		Key:     s,
		Created: time.Now().Unix(),
		Expire:  time.Now().Add(time.Second * time.Duration(api.cacheTTL)).Unix(),
		Body:    resp,
	}
}

// isCachable checks if the passed request is cachable - only GET requests are cachable right now
func isCachable(req *http.Request) bool {
	return req.Method == http.MethodGet
}

// RemoveFromCache removes a single element from the cache
func (api *API) RemoveFromCache(s string) {
	api.muCache.Lock()
	delete(api.cache, s)
	api.cacheGen++
	api.muCache.Unlock()
}

// PruneCache removes all entries from the response cache.
func (api *API) PruneCache() {
	api.muCache.Lock()
	api.cache = make(map[string]*responseCache)
	api.cacheGen++
	api.muCache.Unlock()
}

// BuildCacheKey generates a SHA256 hash from the request URL to use as a cache key.
func BuildCacheKey(req *http.Request) string {
	return BuildSHA256(req.URL.String())
}
