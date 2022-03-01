package myrasec

import (
	"fmt"
	"strconv"

	"github.com/Myra-Security-GmbH/myrasec-go/pkg/types"
)

type Certificate struct {
	ID           int             `json:"id,omitempty"`
	Created      *types.DateTime `json:"created,omitempty"`
	Modified     *types.DateTime `json:"modified,omitempty"`
	Subject      string          `json:"subject"`
	Algorithm    string          `json:"algorithm"`
	ValidFrom    *types.DateTime `json:"validFrom"`
	ValidTo      *types.DateTime `json:"validTo"`
	Fingerprint  string          `json:"finterprint"`
	SerialNumber string          `json:"serialNumber"`
	Cert         string          `json:"cert,omitempty"`
}

type SSLCertificate struct {
	*Certificate
	SubjectAlternatives []string          `json:"subjectAlternatives"`
	Intermediates       []SSLIntermediate `json:"intermediates,omitempty"`
	Wildcard            bool              `json:"wildcard"`
	ExtendedValidation  bool              `json:"extendedValidation"`
	Subdomains          []string          `json:"subdomains,omitempty"`
	Key                 string            `json:"key,omitempty"`
	CertRefreshForced   bool              `json:"certRefreshForced,omitempty"`
	CertToRefresh       int               `json:"certToRefresh,omitempty"`
}

type SSLIntermediate struct {
	*Certificate
	Issuer string `json:"issuer"`
}

//
// ListSSLCertificates returns a slice containing all visible SSL certificates for a domain
//
func (api *API) ListSSLCertificates(domainName string, params map[string]string) ([]SSLCertificate, error) {
	if _, ok := methods["listSSLCertificates"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "listSSLCertificates")
	}

	page := 1
	var err error
	if pageParam, ok := params[ParamPage]; ok {
		delete(params, ParamPage)
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			page = 1
		}
	}

	definition := methods["listSSLCertificates"]
	definition.Action = fmt.Sprintf(definition.Action, domainName, page)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	var records []SSLCertificate
	for _, v := range *result.(*[]SSLCertificate) {
		records = append(records, v)
	}

	return records, nil
}

//
// CreateSSLCertificate creates a new SSL certificates on the passed domain (name) using the MYRA API
//
func (api *API) CreateSSLCertificate(cert *SSLCertificate, domainName string) (*SSLCertificate, error) {
	if _, ok := methods["createSSLCertificate"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "createSSLCertificate")
	}

	definition := methods["createSSLCertificate"]
	definition.Action = fmt.Sprintf(definition.Action, domainName)

	result, err := api.call(definition, cert)
	if err != nil {
		return nil, err
	}
	return result.(*SSLCertificate), nil
}

//
// UpdateSSLCertificate updates the passed SSL certificate using the MYRA API
//
func (api *API) UpdateSSLCertificate(cert *SSLCertificate, domainName string) (*SSLCertificate, error) {
	if _, ok := methods["updateSSLCertificate"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "updateSSLCertificate")
	}

	definition := methods["updateSSLCertificate"]
	definition.Action = fmt.Sprintf(definition.Action, domainName)

	result, err := api.call(definition, cert)
	if err != nil {
		return nil, err
	}
	return result.(*SSLCertificate), nil
}

//
// DeleteSSLCertificate "deletes" the passed SSL certificate by removing the assigned subdomains from the certificate using the MYRA API
//
func (api *API) DeleteSSLCertificate(cert *SSLCertificate, domainName string) (*SSLCertificate, error) {
	if _, ok := methods["deleteSSLCertificate"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "deleteSSLCertificate")
	}

	definition := methods["deleteSSLCertificate"]
	definition.Action = fmt.Sprintf(definition.Action, domainName)

	cert.Subdomains = []string{}

	result, err := api.call(definition, cert)
	if err != nil {
		return nil, err
	}
	return result.(*SSLCertificate), nil
}
