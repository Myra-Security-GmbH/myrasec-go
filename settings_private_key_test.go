package myrasec

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestSettingsPrivateKeyMarkerRoundTrip(t *testing.T) {
	for _, tag := range []bool{false, true} {
		body := `{"myra_ssl_certificate_key":["<redacted>"],"only_https":true}`
		decode := decodeSettingsResponse
		if tag {
			body = `{"settings":` + body + `}`
			decode = decodeTagSettingsResponse
		}
		decoded, err := decode(&http.Response{Body: io.NopCloser(strings.NewReader(body))}, APIMethod{})
		if err != nil {
			t.Fatal(err)
		}
		settings := decoded.(*Settings)
		if !reflect.DeepEqual(settings.MyraSSLCertificateKey, []string{RedactedSettingsPrivateKey}) {
			t.Fatal("configured key marker was lost")
		}
		encoded, err := json.Marshal(settings)
		if err != nil {
			t.Fatal(err)
		}
		var payload map[string]any
		if err := json.Unmarshal(encoded, &payload); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(payload["myra_ssl_certificate_key"], []any{RedactedSettingsPrivateKey}) {
			t.Fatal("update payload must preserve the marker")
		}
	}
}
