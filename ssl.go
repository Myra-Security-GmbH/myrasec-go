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
	ID           int             `json:"id,omitempty" jsonschema:"ID is an unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete a Certificate it is necessary to add this attribute to your object."`
	Created      *types.DateTime `json:"created,omitempty" jsonschema:"Created is a date type attribute with an ISO 8601 format. It will be created by the server after creating a new Certificate object. This value is only informational so it is not necessary to add this an attribute to any API call."`
	Modified     *types.DateTime `json:"modified,omitempty" jsonschema:"Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates and deletions. This value is always a date type with an ISO 8601 format."`
	Subject      string          `json:"subject" jsonschema:"Shows the subject of the uploaded certificate."`
	Algorithm    string          `json:"algorithm" jsonschema:"Contains the signature algorithm."`
	ValidFrom    *types.DateTime `json:"validFrom" jsonschema:"Time when the certificate starts to be valid. This property is a date type with an ISO 8601 format."`
	ValidTo      *types.DateTime `json:"validTo" jsonschema:"Time when the certificate expires. This property is a date type with an ISO 8601 format."`
	Fingerprint  string          `json:"fingerprint" jsonschema:"Fingerprint of the certificate."`
	SerialNumber string          `json:"serialNumber" jsonschema:"Serial number of the certificate."`
	Cert         string          `json:"cert,omitempty" jsonschema:"Cert contains the certificate."`
}

// SSLCertificate struct ...
type SSLCertificate struct {
	*Certificate
	SubjectAlternatives  []string          `json:"subjectAlternatives" jsonschema:"Contains a list of subdomains which can be validated using this certificate. This list also contains the CN of the subject."`
	Intermediates        []SSLIntermediate `json:"intermediates,omitempty" jsonschema:"Contains a list of intermediate certificates to be used in order to generate a chain of trust. The intermediates are filtered and sorted based on subject / issuer relationship. Uploading a partial or a completely different chain will result in an empty list."`
	Wildcard             bool              `json:"wildcard" jsonschema:"This property shows whether the certificate is valid for multiple subdomains of a domain. The certificate needs to have a *.domain.tld subject to return true."`
	ExtendedValidation   bool              `json:"extendedValidation" jsonschema:"True if the browser handles the certificate as extended validation. We use the OIDs from Google Chrome™ to measure the extended validation level."`
	Subdomains           []string          `json:"subdomains,omitempty" jsonschema:"A list of subdomains assigned to this certificate."`
	Key                  string            `json:"key,omitempty" jsonschema:"The unencrypted private key that matches your certificate."`
	CertRefreshForced    bool              `json:"certRefreshForced" jsonschema:"Every time a certificate is refreshed with another non-matching certificate the operation is interrupted with an error. Setting certRefreshForced will ignore such errors and refresh the certificate anyway. Please use it only, if you are sure you can ignore an error when refreshing a certificate."`
	CertToRefresh        int               `json:"certToRefresh,omitempty" jsonschema:"This property allows you to update an already existing certificate with a new one without changing IP addresses, the value has to be the ID of the cert that should be refreshed."`
	SslConfigurationName string            `json:"sslConfigurationName,omitempty" jsonschema:"This property allows you to set a specific ssl configuration. default Myra-Global-TLS-Default, valid values are Myra-Global-TLS-Default, 2023-mozilla-intermediate, 2023-mozilla-modern."`
	Managed              bool              `json:"managed" jsonschema:"Indicates wether this certificate is managed by Myra or not."`
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

	return *result.(*[]SSLCertificate), nil
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
