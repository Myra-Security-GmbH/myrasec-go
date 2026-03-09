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

// Certificate serves as the base structure for SSL certificate data.
// It contains the raw PEM data and extracted metadata like validity dates and fingerprints.
type Certificate struct {
	// ID is the unique identifier for the certificate.
	// This value is server-generated and required for update and delete operations.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the certificate. Server-generated; required for updates and deletes, but ignored during creation."`

	// Created indicates when the certificate object was added to the system.
	// This is a server-managed, read-only value in ISO 8601 format.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified serves as a version identifier for optimistic locking.
	// It records the last update time in ISO 8601 format. This field is required
	// for update and delete operations to ensure data consistency.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Required for updates and deletes to ensure data consistency (optimistic locking)."`

	// Subject is the Common Name (CN) or subject extracted from the certificate.
	// Read-only; automatically parsed from the uploaded PEM data.
	Subject string `json:"subject" jsonschema:"The Common Name (CN) extracted from the certificate. Read-only; automatically populated by the server."`

	// Algorithm specifies the signature algorithm used (e.g., SHA256withRSA).
	// Read-only; automatically parsed from the uploaded PEM data.
	Algorithm string `json:"algorithm" jsonschema:"The signature algorithm (e.g., 'SHA256withRSA'). Read-only; automatically populated by the server."`

	// ValidFrom indicates the start date of the certificate's validity period.
	// Read-only; automatically parsed from the uploaded PEM data.
	ValidFrom *types.DateTime `json:"validFrom" jsonschema:"The timestamp (ISO 8601) when the certificate becomes valid. Read-only; automatically populated by the server."`

	// ValidTo indicates the expiration date of the certificate.
	// Read-only; automatically parsed from the uploaded PEM data.
	ValidTo *types.DateTime `json:"validTo" jsonschema:"The timestamp (ISO 8601) when the certificate expires. Read-only; automatically populated by the server."`

	// Fingerprint is the unique hash (SHA1/SHA256) of the certificate.
	// Read-only; automatically parsed from the uploaded PEM data.
	Fingerprint string `json:"fingerprint" jsonschema:"The unique fingerprint hash of the certificate. Read-only; automatically populated by the server."`

	// SerialNumber is the serial number assigned by the CA.
	// Read-only; automatically parsed from the uploaded PEM data.
	SerialNumber string `json:"serialNumber" jsonschema:"The serial number assigned by the Certificate Authority. Read-only; automatically populated by the server."`

	// Cert contains the raw public certificate data in PEM format.
	// This is the primary input field for uploading a certificate.
	Cert string `json:"cert,omitempty" jsonschema:"The raw public certificate data in PEM format (-----BEGIN CERTIFICATE-----). Required when uploading a new certificate."`
}

// SSLCertificate represents a full SSL configuration object, including the private key
// and assignments to specific subdomains. It embeds the base Certificate metadata.
type SSLCertificate struct {
	*Certificate

	// SubjectAlternatives lists all SANs (Subject Alternative Names) covered by this certificate.
	// Read-only; automatically parsed from the uploaded PEM data.
	SubjectAlternatives []string `json:"subjectAlternatives" jsonschema:"List of Subject Alternative Names (SANs) and the CN covered by this certificate. Read-only; automatically populated by the server."`

	// Intermediates contains the chain of intermediate certificates.
	// The system automatically sorts and filters these based on the uploaded certificate.
	Intermediates []SSLIntermediate `json:"intermediates,omitempty" jsonschema:"List of intermediate certificates required for the chain of trust. Read-only; automatically generated/sorted by the server upon upload."`

	// Wildcard indicates if the certificate is a wildcard certificate (*.domain.tld).
	// Read-only; automatically determined from the subject.
	Wildcard bool `json:"wildcard" jsonschema:"True if the certificate is a wildcard certificate (Subject starts with '*'). Read-only."`

	// ExtendedValidation indicates if the certificate has EV status (green bar).
	// Read-only; determined via OID matching (e.g., Google Chrome™ standards).
	ExtendedValidation bool `json:"extendedValidation" jsonschema:"True if the certificate is detected as Extended Validation (EV). Read-only."`

	// Subdomains is a list of FQDNs in the Myra system assigned to this certificate.
	Subdomains []string `json:"subdomains,omitempty" jsonschema:"List of subdomains (FQDNs) explicitly assigned to use this certificate."`

	// Key is the private key associated with the certificate in PEM format.
	// Required for new uploads. Write-only (usually not returned in GET requests for security).
	Key string `json:"key,omitempty" jsonschema:"The unencrypted private key in PEM format (-----BEGIN RSA PRIVATE KEY-----). Required when uploading a new certificate."`

	// CertRefreshForced allows overwriting an existing certificate even if the new one matches differently.
	// Use with caution to prevent accidental interruptions.
	CertRefreshForced bool `json:"certRefreshForced" jsonschema:"Safety override flag. If true, forces the certificate refresh even if validation errors occur (e.g., non-matching subjects). Use with caution."`

	// CertToRefresh is the ID of an existing certificate to replace/rotate.
	// Used to update a certificate without changing assigned IP addresses.
	CertToRefresh int `json:"certToRefresh,omitempty" jsonschema:"The ID of an existing certificate object to replace. Use this to rotate certificates while preserving IP assignments."`

	// SslConfigurationName specifies the TLS protocol and cipher suite profile.
	// Valid values: 'Myra-Global-TLS-Default', '2023-mozilla-intermediate', '2023-mozilla-modern'.
	SslConfigurationName string `json:"sslConfigurationName,omitempty" jsonschema:"The TLS configuration profile. Valid values: 'Myra-Global-TLS-Default', '2023-mozilla-intermediate', '2023-mozilla-modern'."`

	// Managed indicates if the certificate is automatically managed/renewed by the Myra platform (e.g., Let's Encrypt).
	// Read-only.
	Managed bool `json:"managed" jsonschema:"Indicates if this certificate is automatically managed and renewed by the Myra platform. Read-only."`
}

// SSLIntermediate represents an intermediate CA certificate in the chain.
type SSLIntermediate struct {
	*Certificate

	// Issuer is the name of the entity that signed this intermediate certificate.
	Issuer string `json:"issuer" jsonschema:"The name of the Issuer (CA) that signed this certificate. Read-only."`
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

	res, ok := result.(*SSLCertificate)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
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

	res, ok := result.(*[]SSLCertificate)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
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
	res, ok := result.(*SSLCertificate)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
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
	res, ok := result.(*SSLCertificate)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
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
	res, ok := result.(*SSLCertificate)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}
