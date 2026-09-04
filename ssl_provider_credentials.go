package myrasec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// getSSLProviderCredentialsMethods returns SSL provider credentials related API calls
func getSSLProviderCredentialsMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"listSSLProviderCredentials": {
			Name:   "listSSLProviderCredentials",
			Action: "ssl/providers",
			Method: http.MethodGet,
			Result: []SSLProviderCredentials{},
		},
		"getSSLProviderCredentials": {
			Name:               "getSSLProviderCredentials",
			Action:             "ssl/providers/%d",
			Method:             http.MethodGet,
			Result:             SSLProviderCredentials{},
			ResponseDecodeFunc: decodeSingleElementResponse,
		},
		"createSSLProviderCredentials": {
			Name:   "createSSLProviderCredentials",
			Action: "ssl/providers",
			Method: http.MethodPost,
			Result: SSLProviderCredentials{},
		},
		"updateSSLProviderCredentials": {
			Name:   "updateSSLProviderCredentials",
			Action: "ssl/providers/%d",
			Method: http.MethodPut,
			Result: SSLProviderCredentials{},
		},
		"deleteSSLProviderCredentials": {
			Name:   "deleteSSLProviderCredentials",
			Action: "ssl/providers/%d",
			Method: http.MethodDelete,
			Result: SSLProviderCredentials{},
		},
		"listSSLProviderCertificates": {
			Name:   "listSSLProviderCertificates",
			Action: "ssl/providers/%d/certificates",
			Method: http.MethodGet,
			Result: []SSLCertificateSummary{},
		},
	}
}

// SSLProviderCredentials represents the credentials of an enterprise CA account (Sectigo or
// D-Trust) that Myra Managed Certificates are requested with. A managed certificate request
// references them through its SSLProviderCredentialsID. The endpoints require the
// "Myra-Certificate" organization feature and a root group member or organization admin;
// every other caller receives HTTP 403.
type SSLProviderCredentials struct {
	// ID is the unique identifier for the provider credentials.
	// This value is server-generated and required for update and delete operations.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the provider credentials. Server-generated; required for updates and deletes, but ignored during creation."`

	// Created indicates when the credentials were created.
	// This is a server-managed, read-only value in ISO 8601 format.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified serves as a version identifier for optimistic locking.
	// It records the last update time in ISO 8601 format. This field is required
	// for update and delete operations to ensure data consistency.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Required for updates and deletes to ensure data consistency (optimistic locking)."`

	// Name is an arbitrary, user-defined label for the credentials. Required.
	Name string `json:"name" jsonschema:"An arbitrary, user-defined label for the credentials. Required."`

	// Provider is the certificate provider the credentials belong to.
	// Valid values: 'SECTIGO', 'DTRUST'. Required.
	Provider string `json:"provider" jsonschema:"The certificate provider the credentials belong to. Valid values: 'SECTIGO', 'DTRUST'. Required."`

	// Cert is the public certificate of the ACME account key pair in PEM format.
	// Leave Cert and PrivateKey empty on creation to let the server generate a pair.
	// An empty value on update keeps the stored certificate.
	Cert string `json:"cert,omitempty" jsonschema:"The public certificate of the ACME account key pair (PEM). Leave cert and privateKey empty on creation to let the server generate a pair. Empty on update keeps the stored value."`

	// PrivateKey is the private key of the ACME account key pair in PEM format.
	// Write-only: it is never returned to customer administrators. A value sent on
	// update by a customer administrator is discarded. Contact support to replace the pair.
	PrivateKey string `json:"privateKey,omitempty" jsonschema:"The private key of the ACME account key pair (PEM). Write-only; never returned and not replaceable on update by customer administrators."`

	// Email is the contact email address of the provider account.
	// Optional for Sectigo. Must be empty for D-Trust.
	Email string `json:"email,omitempty" jsonschema:"The contact email address of the provider account. Optional for 'SECTIGO', must be empty for 'DTRUST'."`

	// Endpoint is the ACME directory URL of the provider product, e.g. 'https://acme.sectigo.com/v2/OV'.
	// Required for Sectigo, since region and validation product are encoded in the URL.
	// Optional for D-Trust, which falls back to the platform default.
	Endpoint string `json:"endpoint,omitempty" jsonschema:"The ACME directory URL of the provider product, e.g. 'https://acme.sectigo.com/v2/OV'. Required for 'SECTIGO', optional for 'DTRUST'."`

	// EABKid is the External Account Binding key identifier issued by the provider. Required.
	EABKid string `json:"eabKid,omitempty" jsonschema:"The External Account Binding (EAB) key identifier issued by the provider. Required."`

	// EABHmac is the External Account Binding HMAC key issued by the provider. Required on creation.
	// Write-only: it is never returned in responses. An empty value on update keeps the stored key.
	EABHmac string `json:"eabHmac,omitempty" jsonschema:"The External Account Binding (EAB) HMAC key issued by the provider. Required on creation. Write-only; empty on update keeps the stored value."`

	// Comment is an optional free-text note.
	Comment string `json:"comment,omitempty" jsonschema:"An optional free-text note."`
}

// SSLCertificateSummary is the compact, read-only representation of a certificate issued
// for a managed certificate request. It omits the PEM data, the private key and the
// intermediates of the full SSLCertificate.
type SSLCertificateSummary struct {
	// ID is the unique identifier for the certificate.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the certificate. Read-only."`

	// Created indicates when the certificate object was added to the system. Read-only.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Read-only."`

	// Modified records the last update time. Read-only.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Read-only."`

	// Subject is the Common Name (CN) or subject of the certificate.
	Subject string `json:"subject" jsonschema:"The Common Name (CN) of the certificate. Read-only."`

	// SubjectAlternatives lists all SANs (Subject Alternative Names) covered by this certificate.
	SubjectAlternatives []string `json:"subjectAlternatives" jsonschema:"List of Subject Alternative Names (SANs) covered by this certificate. Read-only."`

	// Algorithm specifies the signature algorithm used (e.g., SHA256withRSA).
	Algorithm string `json:"algorithm" jsonschema:"The signature algorithm (e.g., 'SHA256withRSA'). Read-only."`

	// ValidFrom indicates the start date of the certificate's validity period.
	ValidFrom *types.DateTime `json:"validFrom" jsonschema:"The timestamp (ISO 8601) when the certificate becomes valid. Read-only."`

	// ValidTo indicates the expiration date of the certificate.
	ValidTo *types.DateTime `json:"validTo" jsonschema:"The timestamp (ISO 8601) when the certificate expires. Read-only."`

	// Fingerprint is the unique hash of the certificate.
	Fingerprint string `json:"fingerprint" jsonschema:"The unique fingerprint hash of the certificate. Read-only."`

	// SerialNumber is the serial number assigned by the CA.
	SerialNumber string `json:"serialNumber" jsonschema:"The serial number assigned by the Certificate Authority. Read-only."`

	// Wildcard indicates if the certificate is a wildcard certificate (*.domain.tld).
	Wildcard bool `json:"wildcard" jsonschema:"True if the certificate is a wildcard certificate. Read-only."`

	// ExtendedValidation indicates if the certificate has EV status.
	ExtendedValidation bool `json:"extendedValidation" jsonschema:"True if the certificate is detected as Extended Validation (EV). Read-only."`

	// Managed indicates if the certificate is automatically managed and renewed by the Myra platform.
	Managed bool `json:"managed" jsonschema:"True if the certificate is managed and renewed by the Myra platform. Read-only."`

	// Multidomain indicates if the certificate spans more than one domain.
	Multidomain bool `json:"multidomain" jsonschema:"True if the certificate spans more than one domain. Read-only."`

	// SslConfigurationName is the TLS profile applied to the certificate.
	SslConfigurationName string `json:"sslConfigurationName,omitempty" jsonschema:"The TLS configuration profile applied to the certificate. Read-only."`

	// RequestID is the ID of the managed certificate request the certificate was issued for.
	RequestID int `json:"requestId,omitempty" jsonschema:"The ID of the managed certificate request the certificate was issued for. Read-only."`

	// DomainID is the ID of the domain the certificate belongs to.
	DomainID int `json:"domainId,omitempty" jsonschema:"The ID of the domain the certificate belongs to. Read-only."`

	// Subdomains is a list of FQDNs in the Myra system assigned to this certificate.
	Subdomains []string `json:"subdomains,omitempty" jsonschema:"List of subdomains (FQDNs) assigned to this certificate. Read-only."`
}

// ListSSLProviderCredentialsContext returns a slice containing all SSL provider credentials of the organization.
// The secret fields (EABHmac, PrivateKey) are never returned.
func (api *API) ListSSLProviderCredentialsContext(ctx context.Context, params map[string]string) ([]SSLProviderCredentials, error) {
	if _, ok := api.methods["listSSLProviderCredentials"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listSSLProviderCredentials")
	}

	definition := api.methods["listSSLProviderCredentials"]

	result, err := api.call(ctx, definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]SSLProviderCredentials)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// ListSSLProviderCredentials is equivalent to ListSSLProviderCredentialsContext with context.Background().
//
// Deprecated: use ListSSLProviderCredentialsContext.
func (api *API) ListSSLProviderCredentials(params map[string]string) ([]SSLProviderCredentials, error) {
	return api.ListSSLProviderCredentialsContext(context.Background(), params)
}

// GetSSLProviderCredentialsContext returns a single SSL provider credentials object with/for the given identifier.
// The secret fields (EABHmac, PrivateKey) are never returned.
func (api *API) GetSSLProviderCredentialsContext(ctx context.Context, id int) (*SSLProviderCredentials, error) {
	if _, ok := api.methods["getSSLProviderCredentials"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getSSLProviderCredentials")
	}

	definition := api.methods["getSSLProviderCredentials"]
	definition.Action = fmt.Sprintf(definition.Action, id)

	result, err := api.call(ctx, definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	res, ok := result.(*SSLProviderCredentials)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// GetSSLProviderCredentials is equivalent to GetSSLProviderCredentialsContext with context.Background().
//
// Deprecated: use GetSSLProviderCredentialsContext.
func (api *API) GetSSLProviderCredentials(id int) (*SSLProviderCredentials, error) {
	return api.GetSSLProviderCredentialsContext(context.Background(), id)
}

// CreateSSLProviderCredentialsContext creates new SSL provider credentials using the MYRA API
func (api *API) CreateSSLProviderCredentialsContext(ctx context.Context, credentials *SSLProviderCredentials) (*SSLProviderCredentials, error) {
	if _, ok := api.methods["createSSLProviderCredentials"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createSSLProviderCredentials")
	}

	definition := api.methods["createSSLProviderCredentials"]

	result, err := api.call(ctx, definition, credentials)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*SSLProviderCredentials)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// CreateSSLProviderCredentials is equivalent to CreateSSLProviderCredentialsContext with context.Background().
//
// Deprecated: use CreateSSLProviderCredentialsContext.
func (api *API) CreateSSLProviderCredentials(credentials *SSLProviderCredentials) (*SSLProviderCredentials, error) {
	return api.CreateSSLProviderCredentialsContext(context.Background(), credentials)
}

// UpdateSSLProviderCredentialsContext updates the passed SSL provider credentials using the MYRA API.
// Empty Cert and EABHmac values keep the stored ones, so a partial update never clears a secret.
func (api *API) UpdateSSLProviderCredentialsContext(ctx context.Context, credentials *SSLProviderCredentials) (*SSLProviderCredentials, error) {
	if _, ok := api.methods["updateSSLProviderCredentials"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateSSLProviderCredentials")
	}

	definition := api.methods["updateSSLProviderCredentials"]
	definition.Action = fmt.Sprintf(definition.Action, credentials.ID)

	result, err := api.call(ctx, definition, credentials)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*SSLProviderCredentials)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// UpdateSSLProviderCredentials is equivalent to UpdateSSLProviderCredentialsContext with context.Background().
//
// Deprecated: use UpdateSSLProviderCredentialsContext.
func (api *API) UpdateSSLProviderCredentials(credentials *SSLProviderCredentials) (*SSLProviderCredentials, error) {
	return api.UpdateSSLProviderCredentialsContext(context.Background(), credentials)
}

// DeleteSSLProviderCredentialsContext deletes the passed SSL provider credentials using the MYRA API.
// Managed certificate requests referencing the credentials keep existing, their
// SSLProviderCredentialsID is cleared.
func (api *API) DeleteSSLProviderCredentialsContext(ctx context.Context, credentials *SSLProviderCredentials) (*SSLProviderCredentials, error) {
	if _, ok := api.methods["deleteSSLProviderCredentials"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteSSLProviderCredentials")
	}

	definition := api.methods["deleteSSLProviderCredentials"]
	definition.Action = fmt.Sprintf(definition.Action, credentials.ID)

	_, err := api.call(ctx, definition, credentials)
	if err != nil {
		return nil, err
	}

	return credentials, nil
}

// DeleteSSLProviderCredentials is equivalent to DeleteSSLProviderCredentialsContext with context.Background().
//
// Deprecated: use DeleteSSLProviderCredentialsContext.
func (api *API) DeleteSSLProviderCredentials(credentials *SSLProviderCredentials) (*SSLProviderCredentials, error) {
	return api.DeleteSSLProviderCredentialsContext(context.Background(), credentials)
}

// ListSSLProviderCertificatesContext returns a slice containing the certificates issued with the SSL
// provider credentials with the given identifier.
func (api *API) ListSSLProviderCertificatesContext(ctx context.Context, credentialsId int, params map[string]string) ([]SSLCertificateSummary, error) {
	if _, ok := api.methods["listSSLProviderCertificates"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listSSLProviderCertificates")
	}

	definition := api.methods["listSSLProviderCertificates"]
	definition.Action = fmt.Sprintf(definition.Action, credentialsId)

	result, err := api.call(ctx, definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]SSLCertificateSummary)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// ListSSLProviderCertificates is equivalent to ListSSLProviderCertificatesContext with context.Background().
//
// Deprecated: use ListSSLProviderCertificatesContext.
func (api *API) ListSSLProviderCertificates(credentialsId int, params map[string]string) ([]SSLCertificateSummary, error) {
	return api.ListSSLProviderCertificatesContext(context.Background(), credentialsId, params)
}
