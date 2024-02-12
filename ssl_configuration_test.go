package myrasec

import "testing"

func TestListSslConfigurations(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/ssl-configurations",
			`{"error": false, "pageSize": 10, "page": 1, "count": 1, "data": [
				{
					"id": 1,
					"name": "Myra-Global-TLS-Default",
					"ciphers": "",
					"protocols": ""
				},
				{
					"id": 2,
					"name": "2023-mozilla-intermediate",
					"ciphers": "TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384:DHE-RSA-CHACHA20-POLY1305",
					"protocols": "TLSv1.2,TLSv1.3"
				},
				{
					"id": 3,
					"name": "2023-mozilla-modern",
					"ciphers": "TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256",
					"protocols": "TLSv1.3"
				}
			]}`,
			methods["listSSLConfigurations"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	config, err := api.ListSslConfigurations()
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if len(config) != 3 {
		t.Errorf("Expected to get [%d] ssl configurations but got [%d]", 3, len(config))
	}
}
