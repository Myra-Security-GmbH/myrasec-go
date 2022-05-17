package myrasec

import (
	"net/http"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

//
// responseCache ...
//
type responseCache struct {
	Key     string
	Created int64
	Expire  int64
	Request *http.Request
	Body    interface{}
}

//
// isExpired checks if the cached response is expired
//
func (c *responseCache) isExpired() bool {
	return c.Expire < time.Now().Unix()
}

//
// inCache checks the cache if the response for the passed request is stored in the cache.
//
func (api *API) inCache(req *http.Request) bool {
	s := BuildSHA256(req.URL.String())

	mutex.Lock()
	c, ok := api.cache[s]
	mutex.Unlock()

	if !ok {
		return false
	}

	// if ttl is expired - remove from cache and return false
	if c.isExpired() {
		api.removeFromCache(s)
		return false
	}
	// if the body is nil - return false as we do not have any response cached to return
	if c.Body == nil {
		return false
	}

	return true
}

//
// fromCache loads the response from the cache (if it is cached)
//
func (api *API) fromCache(req *http.Request) interface{} {
	if !api.inCache(req) {
		return nil
	}

	s := BuildSHA256(req.URL.String())
	mutex.Lock()
	defer mutex.Unlock()

	if c, ok := api.cache[s]; ok {
		return c.Body
	}

	return nil
}

//
// cacheResponse stores the response body in the cache
//
func (api *API) cacheResponse(req *http.Request, resp interface{}) {
	if !api.caching {
		return
	}

	s := BuildSHA256(req.URL.String())
	mutex.Lock()
	defer mutex.Unlock()

	api.cache[s] = &responseCache{
		Key:     s,
		Created: time.Now().Unix(),
		Expire:  time.Now().Add(time.Second * time.Duration(api.cacheTTL)).Unix(),
		Request: req,
		Body:    resp,
	}
}

//
// isCachable checks if the passed request is cachable - only GET requests are cachable right now
//
func isCachable(req *http.Request) bool {
	return req.Method == http.MethodGet
}

//
// removeFromCache removes a single element from the cache
//
func (api *API) removeFromCache(s string) {
	mutex.Lock()
	delete(api.cache, s)
	mutex.Unlock()
}
