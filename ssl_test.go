package myrasec

import (
	"strings"
	"testing"
)

func TestGetSSLCertificate(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/domain/1/ssl/certificates/1",
			`{"error": false, "pageSize": 10, "page": 1, "count": 1, "data": [
				{"id": 1, "subjectAlternatives": ["example.com","*.example.com"], "intermediates": [], "wildcard": true, "extendedValidation": false, "subdomains": ["example.com", "www.example.com"], "validFrom": "2022-01-01T00:00:00+0200", "validTo": "2022-01-31T00:00:00+0200", "algorithm": "RSA-SHA256", "subject": "CN=example.com", "fingerprint": "11:22:33:44:55:66:77:88:99:00", "serialNumber": "00"}
			]}`,
			methods["getSSLCertificate"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	cert, err := api.GetSSLCertificate(1, 1)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	if cert.ID != 1 {
		t.Errorf("Expected to get SSL Certificate with ID [%d] but got [%d]", 1, cert.ID)
	}

	if strings.Join(cert.SubjectAlternatives, ",") != "example.com,*.example.com" {
		t.Errorf("Expected to get SSL Cert with  [%s] but got [%s]", "example.com,*.example.com", strings.Join(cert.SubjectAlternatives, ","))
	}

	if len(cert.Intermediates) != 0 {
		t.Errorf("Expected to get SSL Cert without Intermediates")
	}

	if cert.Wildcard != true {
		t.Errorf("Expected to get SSL Cert with Wildcard [%t] but got [%t]", true, cert.Wildcard)
	}

	if cert.ExtendedValidation != false {
		t.Errorf("Expected to get SSL Cert with ExtendedValidation [%t] but got [%t]", false, cert.ExtendedValidation)
	}

	if strings.Join(cert.Subdomains, ",") != "example.com,www.example.com" {
		t.Errorf("Expected to get SSL Cert with Subdomains [%s] but got [%s]", "example.com,www.example.com", strings.Join(cert.Subdomains, ","))
	}

	if cert.ValidFrom.Format("2006-01-02") != "2022-01-01" {
		t.Errorf("Expected to get SSL Cert with ValidFrom [%s] but got [%s]", "2022-01-01", cert.ValidFrom.Format("2006-01-02"))
	}

	if cert.ValidTo.Format("2006-01-02") != "2022-01-31" {
		t.Errorf("Expected to get SSL Cert with ValidTo [%s] but got [%s]", "2022-01-31", cert.ValidTo.Format("2006-01-02"))
	}

	if cert.Algorithm != "RSA-SHA256" {
		t.Errorf("Expected to get SSL Cert with Algorithm [%s] but got [%s]", "RSA-SHA256", cert.Algorithm)
	}

	if cert.SerialNumber != "00" {
		t.Errorf("Expected to get SSL Cert with SerialNumber [%s] but got [%s]", "00", cert.SerialNumber)
	}

	if cert.Fingerprint != "11:22:33:44:55:66:77:88:99:00" {
		t.Errorf("Expected to get SSL Cert with Fingerprint [%s] but got [%s]", "11:22:33:44:55:66:77:88:99:00", cert.Fingerprint)
	}

	if cert.Subject != "CN=example.com" {
		t.Errorf("Expected to get SSL Cert with Subject [%s] but got [%s]", "CN=example.com", cert.Subject)
	}

}

func TestListSSLCertificates(t *testing.T) {
	api, err := setupPreCachedAPI([]*TestCache{
		preCacheRequest(
			"https://apiv2.myracloud.com/domain/1/ssl/certificates",
			`{"error": false, "pageSize": 10, "page": 1, "count": 1, "data": [
				{"id": 1, "subjectAlternatives": ["example.com","*.example.com"], "intermediates": [], "wildcard": true, "extendedValidation": false, "subdomains": ["example.com", "www.example.com"], "validFrom": "2022-01-01T00:00:00+0200", "validTo": "2022-01-31T00:00:00+0200", "algorithm": "RSA-SHA256", "subject": "CN=example.com", "fingerprint": "11:22:33:44:55:66:77:88:99:00", "serialNumber": "00"}
			]}`,
			methods["listSSLCertificates"],
		),
	})
	if err != nil {
		t.Error("Unexpected error.")
	}

	certs, err := api.ListSSLCertificates(1, nil)
	if err != nil {
		t.Errorf("Expected not to get an error but got [%s]", err.Error())
	}

	for _, cert := range certs {

		if cert.ID != 1 {
			t.Errorf("Expected to get SSL Certificate with ID [%d] but got [%d]", 1, cert.ID)
		}

		if strings.Join(cert.SubjectAlternatives, ",") != "example.com,*.example.com" {
			t.Errorf("Expected to get SSL Cert with  [%s] but got [%s]", "example.com,*.example.com", strings.Join(cert.SubjectAlternatives, ","))
		}

		if len(cert.Intermediates) != 0 {
			t.Errorf("Expected to get SSL Cert without Intermediates")
		}

		if cert.Wildcard != true {
			t.Errorf("Expected to get SSL Cert with Wildcard [%t] but got [%t]", true, cert.Wildcard)
		}

		if cert.ExtendedValidation != false {
			t.Errorf("Expected to get SSL Cert with ExtendedValidation [%t] but got [%t]", false, cert.ExtendedValidation)
		}

		if strings.Join(cert.Subdomains, ",") != "example.com,www.example.com" {
			t.Errorf("Expected to get SSL Cert with Subdomains [%s] but got [%s]", "example.com,www.example.com", strings.Join(cert.Subdomains, ","))
		}

		if cert.ValidFrom.Format("2006-01-02") != "2022-01-01" {
			t.Errorf("Expected to get SSL Cert with ValidFrom [%s] but got [%s]", "2022-01-01", cert.ValidFrom.Format("2006-01-02"))
		}

		if cert.ValidTo.Format("2006-01-02") != "2022-01-31" {
			t.Errorf("Expected to get SSL Cert with ValidTo [%s] but got [%s]", "2022-01-31", cert.ValidTo.Format("2006-01-02"))
		}

		if cert.Algorithm != "RSA-SHA256" {
			t.Errorf("Expected to get SSL Cert with Algorithm [%s] but got [%s]", "RSA-SHA256", cert.Algorithm)
		}

		if cert.SerialNumber != "00" {
			t.Errorf("Expected to get SSL Cert with SerialNumber [%s] but got [%s]", "00", cert.SerialNumber)
		}

		if cert.Fingerprint != "11:22:33:44:55:66:77:88:99:00" {
			t.Errorf("Expected to get SSL Cert with Fingerprint [%s] but got [%s]", "11:22:33:44:55:66:77:88:99:00", cert.Fingerprint)
		}

		if cert.Subject != "CN=example.com" {
			t.Errorf("Expected to get SSL Cert with Subject [%s] but got [%s]", "CN=example.com", cert.Subject)
		}
	}
}
