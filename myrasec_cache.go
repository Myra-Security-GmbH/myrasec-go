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
	Request *http.Request
	Body    any
}

// isExpired checks if the cached response is expired
func (c *responseCache) isExpired() bool {
	return c.Expire < time.Now().Unix()
}

// inCache checks the cache if the response for the passed request is stored in the cache.
func (api *API) inCache(req *http.Request) bool {
	s := BuildCacheKey(req)

	api.muCache.Lock()
	c, ok := api.cache[s]
	api.muCache.Unlock()

	if !ok {
		return false
	}

	// if ttl is expired - remove from cache and return false
	if c.isExpired() {
		api.RemoveFromCache(s)
		return false
	}
	// if the body is nil - return false as we do not have any response cached to return
	if c.Body == nil {
		return false
	}

	return true
}

// fromCache loads the response from the cache (if it is cached)
func (api *API) fromCache(req *http.Request) any {
	if !api.inCache(req) {
		return nil
	}

	s := BuildCacheKey(req)
	api.muCache.Lock()
	defer api.muCache.Unlock()

	if c, ok := api.cache[s]; ok {
		return c.Body
	}

	return nil
}

// cacheResponse stores the response body in the cache
func (api *API) cacheResponse(req *http.Request, resp any) {
	if !api.caching {
		return
	}

	s := BuildCacheKey(req)
	api.muCache.Lock()
	defer api.muCache.Unlock()

	api.cache[s] = &responseCache{
		Key:     s,
		Created: time.Now().Unix(),
		Expire:  time.Now().Add(time.Second * time.Duration(api.cacheTTL)).Unix(),
		Request: req,
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
	api.muCache.Unlock()
}

// PruneCache removes all entries from the response cache.
func (api *API) PruneCache() {
	api.cache = make(map[string]*responseCache)
}

// BuildCacheKey generates a SHA256 hash from the request URL to use as a cache key.
func BuildCacheKey(req *http.Request) string {
	return BuildSHA256(req.URL.String())
}
