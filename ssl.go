package myrasec

import (
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// getSSLMethods returns SSL certificate related API calls
func getSSLMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"getSSLCertificate": {
			Name:               "getSSLCertificate",
			Action:             "domain/%d/ssl/certificates/%d",
			Method:             http.MethodGet,
			Result:             SSLCertificate{},
			ResponseDecodeFunc: decodeSingleElementResponse,
		},
		"listSSLCertificates": {
			Name:   "listSSLCertificates",
			Action: "domain/%d/ssl/certificates",
			Method: http.MethodGet,
			Result: []SSLCertificate{},
		},
		"createSSLCertificate": {
			Name:   "createSSLCertificate",
			Action: "domain/%d/certificates",
			Method: http.MethodPost,
			Result: SSLCertificate{},
		},
		"updateSSLCertificate": {
			Name:   "updateSSLCertificate",
			Action: "domain/%d/certificates/%d",
			Method: http.MethodPut,
			Result: SSLCertificate{},
		},
		"deleteSSLCertificate": {
			Name:   "deleteSSLCertificate",
			Action: "domain/%d/certificates/%d",
			Method: http.MethodPut,
			Result: SSLCertificate{},
		},
	}
}

// Certificate strict ...
type Certificate struct {
	ID           int             `json:"id,omitempty"`
	Created      *types.DateTime `json:"created,omitempty"`
	Modified     *types.DateTime `json:"modified,omitempty"`
	Subject      string          `json:"subject"`
	Algorithm    string          `json:"algorithm"`
	ValidFrom    *types.DateTime `json:"validFrom"`
	ValidTo      *types.DateTime `json:"validTo"`
	Fingerprint  string          `json:"fingerprint"`
	SerialNumber string          `json:"serialNumber"`
	Cert         string          `json:"cert,omitempty"`
}

// SSLCertificate struct ...
type SSLCertificate struct {
	*Certificate
	SubjectAlternatives  []string          `json:"subjectAlternatives"`
	Intermediates        []SSLIntermediate `json:"intermediates,omitempty"`
	Wildcard             bool              `json:"wildcard"`
	ExtendedValidation   bool              `json:"extendedValidation"`
	Subdomains           []string          `json:"subdomains,omitempty"`
	Key                  string            `json:"key,omitempty"`
	CertRefreshForced    bool              `json:"certRefreshForced"`
	CertToRefresh        int               `json:"certToRefresh,omitempty"`
	SslConfigurationName string            `json:"sslConfigurationName,omitempty"`
}

// SSLIntermediate struct ...
type SSLIntermediate struct {
	*Certificate
	Issuer string `json:"issuer"`
}

// GetSSLCertificate returns a single SSL certificate with/for the given identifier
func (api *API) GetSSLCertificate(domainId int, id int) (*SSLCertificate, error) {
	if _, ok := methods["getSSLCertificate"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getSSLCertificate")
	}

	definition := methods["getSSLCertificate"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, id)

	result, err := api.call(definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	return result.(*SSLCertificate), nil
}

// ListSSLCertificates returns a slice containing all visible SSL certificates for a domain
func (api *API) ListSSLCertificates(domainId int, params map[string]string) ([]SSLCertificate, error) {
	if _, ok := methods["listSSLCertificates"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listSSLCertificates")
	}

	definition := methods["listSSLCertificates"]
	definition.Action = fmt.Sprintf(definition.Action, domainId)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	var records []SSLCertificate
	records = append(records, *result.(*[]SSLCertificate)...)

	return records, nil
}

// CreateSSLCertificate creates a new SSL certificates on the passed domain (ID) using the MYRA API
func (api *API) CreateSSLCertificate(cert *SSLCertificate, domainId int) (*SSLCertificate, error) {
	if _, ok := methods["createSSLCertificate"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createSSLCertificate")
	}

	definition := methods["createSSLCertificate"]
	definition.Action = fmt.Sprintf(definition.Action, domainId)

	result, err := api.call(definition, cert)
	if err != nil {
		return nil, err
	}
	return result.(*SSLCertificate), nil
}

// UpdateSSLCertificate updates the passed SSL certificate using the MYRA API
func (api *API) UpdateSSLCertificate(cert *SSLCertificate, domainId int) (*SSLCertificate, error) {
	if _, ok := methods["updateSSLCertificate"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateSSLCertificate")
	}

	definition := methods["updateSSLCertificate"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, cert.ID)

	result, err := api.call(definition, cert)
	if err != nil {
		return nil, err
	}
	return result.(*SSLCertificate), nil
}

// DeleteSSLCertificate "deletes" the passed SSL certificate by removing the assigned subdomains from the certificate using the MYRA API
func (api *API) DeleteSSLCertificate(cert *SSLCertificate, domainId int) (*SSLCertificate, error) {
	if _, ok := methods["deleteSSLCertificate"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteSSLCertificate")
	}

	definition := methods["deleteSSLCertificate"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, cert.ID)

	cert.Subdomains = []string{}

	result, err := api.call(definition, cert)
	if err != nil {
		return nil, err
	}
	return result.(*SSLCertificate), nil
}
