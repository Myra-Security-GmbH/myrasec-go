package myrasec

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	// maxRedirects is the number of redirects the HTTP client follows before it gives up,
	// the limit the default policy of net/http uses as well.
	maxRedirects = 10

	// maxRetryWait caps the wait before the next attempt, whether a Retry-After header asks
	// for it or the doubling of the configured sleep reaches it.
	maxRetryWait = 60 * time.Second

	// maxBackoffShift stops the shift that doubles the configured sleep from overflowing.
	// The wait itself is capped by maxRetryWait long before the shift matters.
	maxBackoffShift = 6

	// maxDrainBytes bounds how much of a discarded response body is read so the connection
	// can be reused for the next attempt.
	maxDrainBytes = 64 << 10
)

// ErrInsecureEndpoint is returned when a request would go to a plain HTTP endpoint on a host
// other than the local machine and SetInsecureHTTP(true) was not called.
var ErrInsecureEndpoint = errors.New("plain HTTP to a host other than the local machine refused, the API credentials would travel in clear text: use an https endpoint or call SetInsecureHTTP(true)")

// checkRedirect is the redirect policy of the HTTP client. The API does not redirect, so a
// redirect means a misconfigured endpoint or a hostile network. net/http keeps the
// Authorization header and the request signature on a redirect to the same host regardless
// of the scheme, so a redirect from HTTPS to plain HTTP would send the credentials in clear
// text. Such downgrades and redirects to another host are refused, same-host HTTPS redirects
// are followed up to maxRedirects.
func checkRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= maxRedirects {
		return fmt.Errorf("stopped after %d redirects", maxRedirects)
	}

	first := via[0]
	if first.URL.Scheme == "https" && req.URL.Scheme != "https" {
		return fmt.Errorf("refusing redirect from %s to %s: downgrade from HTTPS to plain HTTP", first.URL.Redacted(), req.URL.Redacted())
	}

	if !strings.EqualFold(req.URL.Host, first.URL.Host) {
		return fmt.Errorf("refusing redirect from %s to %s: different host", first.URL.Redacted(), req.URL.Redacted())
	}

	return nil
}

// checkScheme refuses requests that would send the credentials in clear text: plain HTTP is
// only allowed to the local machine or after SetInsecureHTTP(true).
func (api *API) checkScheme(u *url.URL) error {
	switch u.Scheme {
	case "https":
		return nil
	case "http":
		if api.insecureHTTP || isLoopbackHost(u.Hostname()) {
			return nil
		}

		return fmt.Errorf("%w: %s", ErrInsecureEndpoint, u.Redacted())
	default:
		return fmt.Errorf("unsupported URL scheme %q in %s", u.Scheme, u.Redacted())
	}
}

// isLoopbackHost reports whether host names the local machine: "localhost" or a loopback IP.
func isLoopbackHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}

	ip := net.ParseIP(host)

	return ip != nil && ip.IsLoopback()
}

// retryableStatus reports whether a response is worth another attempt. A 429 comes from a rate
// limiter that answers before the request is processed, so every method is repeated. A 503
// with Retry-After is the explicit invitation of a server that did not handle the request and
// is repeated for every method as well. Without the header a 503 can come from a gateway that
// lost the API after passing the request on, or from the API after a partial change, just
// like 500, 502 and 504: a repeated POST would create a duplicate, a repeated PUT or DELETE
// fails on the changed or missing entry. Those responses are retried for GET only.
func retryableStatus(method string, status int, retryAfter bool) bool {
	switch status {
	case http.StatusTooManyRequests:
		return true
	case http.StatusServiceUnavailable:
		return method == http.MethodGet || retryAfter
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusGatewayTimeout:
		return method == http.MethodGet
	default:
		return false
	}
}

// retryDelay returns the wait before the next attempt, attempt being the number of the
// attempt that just failed, starting at 1. A Retry-After header of the response (seconds or
// an HTTP date) is honoured up to maxRetryWait. Without the header the wait is the configured
// retry sleep doubled with every attempt (1s, 2s, 4s, ...) up to maxRetryWait, or the
// configured sleep alone when that is longer. A 429 or 503 without the header and without a
// configured sleep waits one second, repeating such a request at once is pointless.
func (api *API) retryDelay(status int, header http.Header, attempt int) time.Duration {
	if wait, ok := parseRetryAfter(header.Get("Retry-After")); ok {
		return min(wait, maxRetryWait)
	}

	sleep := time.Duration(api.retrySleep) * time.Second
	if sleep > maxRetryWait {
		return sleep
	}

	wait := min(sleep<<min(attempt-1, maxBackoffShift), maxRetryWait)
	if wait > 0 {
		return wait
	}

	if status == http.StatusTooManyRequests || status == http.StatusServiceUnavailable {
		return time.Second
	}

	return 0
}

// limiterWait returns how long the rate limiter would hold a request made at the given time.
func (api *API) limiterWait(at time.Time) time.Duration {
	tokens := api.limiter.TokensAt(at)
	limit := api.limiter.Limit()
	if tokens >= 1 || limit <= 0 {
		return 0
	}

	return time.Duration((1 - tokens) / float64(limit) * float64(time.Second))
}

// parseRetryAfter reads a Retry-After header value, either a number of seconds or an HTTP date.
// HTTP/1 trims the value in net/http, HTTP/2 hands it over verbatim.
func parseRetryAfter(value string) (time.Duration, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false
	}

	if seconds, err := strconv.Atoi(value); err == nil {
		if seconds > int(maxRetryWait/time.Second) {
			return maxRetryWait, true
		}

		return time.Duration(max(seconds, 0)) * time.Second, true
	}

	if at, err := http.ParseTime(value); err == nil {
		return max(time.Until(at), 0), true
	}

	return 0, false
}
