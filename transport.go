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

	// maxRetryAfter caps the wait a Retry-After header can request before the next attempt.
	maxRetryAfter = 60 * time.Second

	// maxBackoffShift caps the exponential growth of the sleep between attempts.
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

// retryableStatus reports whether a response is worth another attempt. 429 and 503 mean the
// API did not process the request, so every method is repeated. 500, 502 and 504 can be
// reported for a request the API has processed already (the gateway in front of the API timed
// out, or the API failed after it stored the change): a repeated POST would create a
// duplicate, a repeated PUT or DELETE fails on the changed or missing entry. Those statuses
// are retried for GET only.
func retryableStatus(method string, status int) bool {
	switch status {
	case http.StatusTooManyRequests, http.StatusServiceUnavailable:
		return true
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusGatewayTimeout:
		return method == http.MethodGet
	default:
		return false
	}
}

// retryDelay returns the wait before the next attempt. A Retry-After header of the response
// (seconds or an HTTP date) is honoured up to maxRetryAfter. Without the header the wait is
// the configured retry sleep doubled with every attempt (1s, 2s, 4s, ...) up to maxRetryAfter,
// or the configured sleep alone when that is longer. A 429 or 503 without the header and
// without a configured sleep waits one second, repeating such a request at once is pointless.
func (api *API) retryDelay(status int, header http.Header, attempt int) time.Duration {
	if wait, ok := parseRetryAfter(header.Get("Retry-After")); ok {
		return min(wait, maxRetryAfter)
	}

	sleep := time.Duration(api.retrySleep) * time.Second
	if sleep > maxRetryAfter {
		return sleep
	}

	wait := min(sleep<<min(max(attempt-1, 0), maxBackoffShift), maxRetryAfter)
	if wait > 0 {
		return wait
	}

	if status == http.StatusTooManyRequests || status == http.StatusServiceUnavailable {
		return time.Second
	}

	return 0
}

// parseRetryAfter reads a Retry-After header value, either a number of seconds or an HTTP date.
func parseRetryAfter(value string) (time.Duration, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false
	}

	if seconds, err := strconv.Atoi(value); err == nil {
		if seconds > int(maxRetryAfter/time.Second) {
			return maxRetryAfter, true
		}

		return time.Duration(max(seconds, 0)) * time.Second, true
	}

	if at, err := http.ParseTime(value); err == nil {
		return max(time.Until(at), 0), true
	}

	return 0, false
}
