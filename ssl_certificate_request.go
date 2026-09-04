package myrasec

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// Certificate providers a managed certificate can be requested from.
// Sectigo and D-Trust additionally require stored SSLProviderCredentials.
const (
	SSLProviderLetsEncrypt = "LETS_ENCRYPT"
	SSLProviderSectigo     = "SECTIGO"
	SSLProviderDTrust      = "DTRUST"
)

// Key algorithms of a managed certificate request.
// RSA4096 and RSA8192 are accepted by Sectigo and D-Trust only.
const (
	SSLCertificateRequestAlgorithmRSA2048  = "RSA2048"
	SSLCertificateRequestAlgorithmRSA4096  = "RSA4096"
	SSLCertificateRequestAlgorithmRSA8192  = "RSA8192"
	SSLCertificateRequestAlgorithmECDSA256 = "ECDSA256"
	SSLCertificateRequestAlgorithmECDSA384 = "ECDSA384"
)

// Signature algorithms of a managed certificate request. Accepted by Sectigo and D-Trust only.
const (
	SSLCertificateRequestSignatureAlgorithmSHA256 = "SHA256"
	SSLCertificateRequestSignatureAlgorithmSHA384 = "SHA384"
	SSLCertificateRequestSignatureAlgorithmSHA512 = "SHA512"
)

// Statuses of a managed certificate request.
const (
	// SSLCertificateRequestStatusOpen means the request has been accepted and the certificate is being issued.
	SSLCertificateRequestStatusOpen = "OPEN"
	// SSLCertificateRequestStatusWaitingForCNAME means the issuance waits for the domain validation to succeed.
	SSLCertificateRequestStatusWaitingForCNAME = "WAITING_FOR_CNAME"
	// SSLCertificateRequestStatusCreated means the certificate has been issued and assigned to the subdomains.
	SSLCertificateRequestStatusCreated = "CREATED"
	// SSLCertificateRequestStatusFailed means the issuance did not succeed. The request is not retried automatically.
	SSLCertificateRequestStatusFailed = "FAILED"
)

// Failure reasons of a managed certificate request in status FAILED.
// CNAME_TIMEOUT and VALIDATION_FAILED can be resolved by the customer (DNS or CNAME issue),
// the others are handled on the Myra or CA side.
const (
	SSLCertificateRequestFailureReasonValidationFailed = "VALIDATION_FAILED"
	SSLCertificateRequestFailureReasonValidationError  = "VALIDATION_ERROR"
	SSLCertificateRequestFailureReasonOrderFailed      = "ORDER_FAILED"
	SSLCertificateRequestFailureReasonCNAMETimeout     = "CNAME_TIMEOUT"
	SSLCertificateRequestFailureReasonCertLoadFailed   = "CERT_LOAD_FAILED"
)

// getSSLCertificateRequestMethods returns managed certificate request related API calls
func getSSLCertificateRequestMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"listSSLCertificateRequests": {
			Name:   "listSSLCertificateRequests",
			Action: "ssl/requests",
			Method: http.MethodGet,
			Result: []SSLCertificateRequest{},
		},
		"getSSLCertificateRequest": {
			Name:               "getSSLCertificateRequest",
			Action:             "ssl/requests/%d",
			Method:             http.MethodGet,
			Result:             SSLCertificateRequest{},
			ResponseDecodeFunc: decodeSingleElementResponse,
		},
		"createSSLCertificateRequest": {
			Name:   "createSSLCertificateRequest",
			Action: "ssl/requests",
			Method: http.MethodPost,
			Result: SSLCertificateRequest{},
		},
		"updateSSLCertificateRequest": {
			Name:   "updateSSLCertificateRequest",
			Action: "ssl/requests/%d",
			Method: http.MethodPut,
			Result: SSLCertificateRequest{},
		},
		"deleteSSLCertificateRequest": {
			Name:   "deleteSSLCertificateRequest",
			Action: "ssl/requests/%d",
			Method: http.MethodDelete,
			Result: SSLCertificateRequest{},
		},
		"updateSSLCertificateRequestConfiguration": {
			Name:   "updateSSLCertificateRequestConfiguration",
			Action: "ssl/requests/%d/ssl-configuration",
			Method: http.MethodPut,
			Result: SSLCertificateRequest{},
		},
		"checkSSLCertificateRequestDomains": {
			Name:   "checkSSLCertificateRequestDomains",
			Action: "ssl/requests/check-domains",
			Method: http.MethodPost,
			Result: sslCertificateRequestDomainCheckResponse{},
		},
	}
}

// SSLCertificateRequest represents a Myra Managed Certificate (MMC) request.
// It describes which certificate is requested from which provider and tracks the
// asynchronous issuance through its Status. The organization needs the
// "Myra-Certificate" feature, otherwise every call answers with HTTP 403.
type SSLCertificateRequest struct {
	// ID is the unique identifier for the certificate request.
	// This value is server-generated and required for update and delete operations.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the certificate request. Server-generated; required for updates and deletes, but ignored during creation."`

	// Created indicates when the request was created.
	// This is a server-managed, read-only value in ISO 8601 format.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified serves as a version identifier for optimistic locking.
	// It records the last update time in ISO 8601 format. This field is required
	// for update and delete operations to ensure data consistency.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Required for updates and deletes to ensure data consistency (optimistic locking)."`

	// Algorithm is the key algorithm of the requested certificate.
	// Valid values: 'RSA2048', 'RSA4096', 'RSA8192', 'ECDSA256', 'ECDSA384'. RSA4096 and RSA8192
	// are accepted by Sectigo and D-Trust only. The algorithm is immutable after creation and
	// has to be sent unchanged on updates.
	Algorithm string `json:"algorithm" jsonschema:"The key algorithm of the requested certificate. Valid values: 'RSA2048', 'RSA4096', 'RSA8192', 'ECDSA256', 'ECDSA384'. Immutable after creation; send it unchanged on updates."`

	// Provider is the certificate provider the certificate is requested from.
	// Valid values: 'LETS_ENCRYPT', 'SECTIGO', 'DTRUST'. Sectigo and D-Trust require
	// SSLProviderCredentialsID to reference stored provider credentials.
	Provider string `json:"provider" jsonschema:"The certificate provider. Valid values: 'LETS_ENCRYPT', 'SECTIGO', 'DTRUST'. Sectigo and D-Trust additionally require sslProviderCredentialsId."`

	// Status reports the progress of the issuance.
	// Valid values: 'OPEN', 'WAITING_FOR_CNAME', 'CREATED', 'FAILED'. Read-only; status
	// transitions happen asynchronously on the server.
	Status string `json:"status,omitempty" jsonschema:"The issuance status. Valid values: 'OPEN', 'WAITING_FOR_CNAME', 'CREATED', 'FAILED'. Read-only."`

	// FailureReason carries the reason code when Status is 'FAILED'.
	// Valid values: 'VALIDATION_FAILED', 'VALIDATION_ERROR', 'ORDER_FAILED', 'CNAME_TIMEOUT', 'CERT_LOAD_FAILED'. Read-only.
	FailureReason string `json:"failureReason,omitempty" jsonschema:"The reason code when the status is 'FAILED'. Valid values: 'VALIDATION_FAILED', 'VALIDATION_ERROR', 'ORDER_FAILED', 'CNAME_TIMEOUT', 'CERT_LOAD_FAILED'. Read-only."`

	// CustomerActionable is true when the FailureReason is one the customer can resolve
	// themselves (a DNS or CNAME issue) rather than a Myra or CA-side failure. Read-only.
	CustomerActionable bool `json:"customerActionable,omitempty" jsonschema:"True if the failure reason is a DNS or CNAME issue the customer can resolve. Read-only."`

	// SubjectAlternativeNames lists the names the certificate has to cover.
	// Every name must resolve to a domain registered in the Myra system. Adding a name
	// that no issued certificate of the request covers triggers a new issuance.
	SubjectAlternativeNames []SSLCertificateRequestSAN `json:"subjectAlternativeNames" jsonschema:"The subject alternative names the certificate has to cover. Every name must belong to a domain registered in the Myra system."`

	// Assignments lists the subdomains the issued certificate is assigned to.
	Assignments []SSLCertificateRequestAssignment `json:"assignments" jsonschema:"The subdomains the issued certificate is assigned to."`

	// MultiDomain is true when the request spans more than one domain. Read-only.
	MultiDomain bool `json:"multiDomain,omitempty" jsonschema:"True if the request spans more than one domain. Read-only."`

	// SSLProviderCredentialsID references the SSLProviderCredentials used for the issuance.
	// Required for Sectigo and D-Trust, ignored for Let's Encrypt. Cleared by the server
	// when the referenced credentials are deleted.
	SSLProviderCredentialsID int `json:"sslProviderCredentialsId,omitempty" jsonschema:"The ID of the SSL provider credentials used for the issuance. Required for 'SECTIGO' and 'DTRUST'."`

	// RenewalInterval is the number of days before the certificate expires at which it is
	// renewed. Zero means the system default. Accepted by Sectigo and D-Trust only.
	RenewalInterval int `json:"renewalInterval,omitempty" jsonschema:"Days before expiry at which the certificate is renewed. 0 means the system default. Accepted by 'SECTIGO' and 'DTRUST' only."`

	// SignatureAlgorithm is the signature algorithm of the requested certificate.
	// Valid values: 'SHA256', 'SHA384', 'SHA512'. Empty means the system default.
	// Accepted by Sectigo and D-Trust only. SHA512 cannot be combined with an ECDSA key algorithm.
	SignatureAlgorithm string `json:"signatureAlgorithm,omitempty" jsonschema:"The signature algorithm. Valid values: 'SHA256', 'SHA384', 'SHA512'. Empty means the system default. Accepted by 'SECTIGO' and 'DTRUST' only. 'SHA512' cannot be combined with an ECDSA key algorithm."`
}

// SSLCertificateRequestSAN is a subject alternative name of a managed certificate request.
type SSLCertificateRequestSAN struct {
	// ID is the unique identifier for the subject alternative name. Server-generated.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the subject alternative name. Server-generated."`

	// Created indicates when the subject alternative name was added. Read-only.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified records the last update time. Read-only.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Server-managed, read-only."`

	// Name is the FQDN the certificate has to cover, e.g. 'www.example.com' or '*.example.com'.
	Name string `json:"name" jsonschema:"The FQDN the certificate has to cover, e.g. 'www.example.com' or '*.example.com'."`
}

// SSLCertificateRequestAssignment assigns the issued certificate of a managed
// certificate request to a subdomain.
type SSLCertificateRequestAssignment struct {
	// ID is the unique identifier for the assignment. Server-generated.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the assignment. Server-generated."`

	// Created indicates when the assignment was added. Read-only.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified records the last update time. Read-only.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Server-managed, read-only."`

	// SubDomainName is the FQDN of the subdomain the certificate is assigned to.
	SubDomainName string `json:"subDomainName" jsonschema:"The FQDN of the subdomain the issued certificate is assigned to."`
}

// SSLCertificateRequestDomainCheck is the name server check result for one domain.
type SSLCertificateRequestDomainCheck struct {
	// Exists is true when the domain has name server records.
	Exists bool `json:"exists" jsonschema:"True if the domain has name server records."`

	// IsMyraNS is true when the domain is served by Myra name servers, in which case
	// the domain validation needs no CNAME record.
	IsMyraNS bool `json:"isMyraNS" jsonschema:"True if the domain is served by Myra name servers. No CNAME record is needed then."`

	// ChallengeName is the record name that has to point to ExpectedCName for the
	// domain validation when the domain is not served by Myra name servers.
	ChallengeName string `json:"challengeName,omitempty" jsonschema:"The record name that has to be created as CNAME for the domain validation."`

	// ExpectedCName is the target the ChallengeName record has to point to.
	ExpectedCName string `json:"expectedCName,omitempty" jsonschema:"The target the challenge CNAME record has to point to."`
}

// SSLCertificateRequestDomainChecks maps a checked domain to its name server check result.
// Domains whose lookup failed are absent from the map.
type SSLCertificateRequestDomainChecks map[string]SSLCertificateRequestDomainCheck

// UnmarshalJSON accepts the empty JSON array the API returns when no domain could be resolved.
func (c *SSLCertificateRequestDomainChecks) UnmarshalJSON(b []byte) error {
	trimmed := bytes.TrimSpace(b)
	if bytes.Equal(trimmed, []byte("[]")) || bytes.Equal(trimmed, []byte("null")) {
		*c = SSLCertificateRequestDomainChecks{}
		return nil
	}

	var checks map[string]SSLCertificateRequestDomainCheck
	if err := json.Unmarshal(b, &checks); err != nil {
		return err
	}

	*c = checks

	return nil
}

// sslCertificateRequestDomainCheckResponse is the envelope of the check-domains call.
type sslCertificateRequestDomainCheckResponse struct {
	Domains SSLCertificateRequestDomainChecks `json:"domains"`
}

// payload returns a copy of the request whose nil collections are replaced by empty ones.
// The API expects subjectAlternativeNames and assignments as JSON arrays.
func (r *SSLCertificateRequest) payload() *SSLCertificateRequest {
	payload := *r
	if payload.SubjectAlternativeNames == nil {
		payload.SubjectAlternativeNames = []SSLCertificateRequestSAN{}
	}
	if payload.Assignments == nil {
		payload.Assignments = []SSLCertificateRequestAssignment{}
	}

	return &payload
}

// ListSSLCertificateRequestsContext returns a slice containing all visible managed certificate requests.
// Without a "status" parameter the API omits requests in status CREATED.
func (api *API) ListSSLCertificateRequestsContext(ctx context.Context, params map[string]string) ([]SSLCertificateRequest, error) {
	if _, ok := api.methods["listSSLCertificateRequests"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listSSLCertificateRequests")
	}

	definition := api.methods["listSSLCertificateRequests"]

	result, err := api.call(ctx, definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]SSLCertificateRequest)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// ListSSLCertificateRequests is equivalent to ListSSLCertificateRequestsContext with context.Background().
//
// Deprecated: use ListSSLCertificateRequestsContext.
func (api *API) ListSSLCertificateRequests(params map[string]string) ([]SSLCertificateRequest, error) {
	return api.ListSSLCertificateRequestsContext(context.Background(), params)
}

// GetSSLCertificateRequestContext returns a single managed certificate request with/for the given identifier
func (api *API) GetSSLCertificateRequestContext(ctx context.Context, id int) (*SSLCertificateRequest, error) {
	if _, ok := api.methods["getSSLCertificateRequest"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getSSLCertificateRequest")
	}

	definition := api.methods["getSSLCertificateRequest"]
	definition.Action = fmt.Sprintf(definition.Action, id)

	result, err := api.call(ctx, definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	res, ok := result.(*SSLCertificateRequest)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// GetSSLCertificateRequest is equivalent to GetSSLCertificateRequestContext with context.Background().
//
// Deprecated: use GetSSLCertificateRequestContext.
func (api *API) GetSSLCertificateRequest(id int) (*SSLCertificateRequest, error) {
	return api.GetSSLCertificateRequestContext(context.Background(), id)
}

// CreateSSLCertificateRequestContext creates a new managed certificate request using the MYRA API.
// The issuance is asynchronous: poll GetSSLCertificateRequest until the Status is CREATED or FAILED.
func (api *API) CreateSSLCertificateRequestContext(ctx context.Context, request *SSLCertificateRequest) (*SSLCertificateRequest, error) {
	if _, ok := api.methods["createSSLCertificateRequest"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createSSLCertificateRequest")
	}

	definition := api.methods["createSSLCertificateRequest"]

	result, err := api.call(ctx, definition, request.payload())
	if err != nil {
		return nil, err
	}

	res, ok := result.(*SSLCertificateRequest)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// CreateSSLCertificateRequest is equivalent to CreateSSLCertificateRequestContext with context.Background().
//
// Deprecated: use CreateSSLCertificateRequestContext.
func (api *API) CreateSSLCertificateRequest(request *SSLCertificateRequest) (*SSLCertificateRequest, error) {
	return api.CreateSSLCertificateRequestContext(context.Background(), request)
}

// UpdateSSLCertificateRequestContext updates the passed managed certificate request using the MYRA API.
// The Algorithm is immutable and has to be sent unchanged.
func (api *API) UpdateSSLCertificateRequestContext(ctx context.Context, request *SSLCertificateRequest) (*SSLCertificateRequest, error) {
	if _, ok := api.methods["updateSSLCertificateRequest"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateSSLCertificateRequest")
	}

	definition := api.methods["updateSSLCertificateRequest"]
	definition.Action = fmt.Sprintf(definition.Action, request.ID)

	result, err := api.call(ctx, definition, request.payload())
	if err != nil {
		return nil, err
	}

	res, ok := result.(*SSLCertificateRequest)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// UpdateSSLCertificateRequest is equivalent to UpdateSSLCertificateRequestContext with context.Background().
//
// Deprecated: use UpdateSSLCertificateRequestContext.
func (api *API) UpdateSSLCertificateRequest(request *SSLCertificateRequest) (*SSLCertificateRequest, error) {
	return api.UpdateSSLCertificateRequestContext(context.Background(), request)
}

// DeleteSSLCertificateRequestContext deletes the passed managed certificate request using the MYRA API.
// The certificates issued for the request are removed as well.
func (api *API) DeleteSSLCertificateRequestContext(ctx context.Context, request *SSLCertificateRequest) (*SSLCertificateRequest, error) {
	if _, ok := api.methods["deleteSSLCertificateRequest"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteSSLCertificateRequest")
	}

	definition := api.methods["deleteSSLCertificateRequest"]
	definition.Action = fmt.Sprintf(definition.Action, request.ID)

	_, err := api.call(ctx, definition, request.payload())
	if err != nil {
		return nil, err
	}

	return request, nil
}

// DeleteSSLCertificateRequest is equivalent to DeleteSSLCertificateRequestContext with context.Background().
//
// Deprecated: use DeleteSSLCertificateRequestContext.
func (api *API) DeleteSSLCertificateRequest(request *SSLCertificateRequest) (*SSLCertificateRequest, error) {
	return api.DeleteSSLCertificateRequestContext(context.Background(), request)
}

// UpdateSSLCertificateRequestConfigurationContext applies the SSL configuration (TLS profile) with the
// passed name to every certificate issued for the managed certificate request.
// Valid names are returned by ListSslConfigurations.
func (api *API) UpdateSSLCertificateRequestConfigurationContext(ctx context.Context, id int, sslConfigurationName string) (*SSLCertificateRequest, error) {
	if _, ok := api.methods["updateSSLCertificateRequestConfiguration"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateSSLCertificateRequestConfiguration")
	}

	definition := api.methods["updateSSLCertificateRequestConfiguration"]
	definition.Action = fmt.Sprintf(definition.Action, id)

	result, err := api.call(ctx, definition, map[string]string{"sslConfigurationName": sslConfigurationName})
	if err != nil {
		return nil, err
	}

	res, ok := result.(*SSLCertificateRequest)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// UpdateSSLCertificateRequestConfiguration is equivalent to UpdateSSLCertificateRequestConfigurationContext with context.Background().
//
// Deprecated: use UpdateSSLCertificateRequestConfigurationContext.
func (api *API) UpdateSSLCertificateRequestConfiguration(id int, sslConfigurationName string) (*SSLCertificateRequest, error) {
	return api.UpdateSSLCertificateRequestConfigurationContext(context.Background(), id, sslConfigurationName)
}

// CheckSSLCertificateRequestDomainsContext runs a live name server check for the passed domains and
// returns the result per domain. Use it before creating a request to find out whether the
// domain validation needs a CNAME record. At most 99 domains can be checked per call.
func (api *API) CheckSSLCertificateRequestDomainsContext(ctx context.Context, domains []string) (SSLCertificateRequestDomainChecks, error) {
	if _, ok := api.methods["checkSSLCertificateRequestDomains"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "checkSSLCertificateRequestDomains")
	}

	if len(domains) == 0 {
		return nil, errors.New("no domains passed to check")
	}

	definition := api.methods["checkSSLCertificateRequestDomains"]

	result, err := api.call(ctx, definition, map[string][]string{"domains": domains})
	if err != nil {
		return nil, err
	}

	res, ok := result.(*sslCertificateRequestDomainCheckResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res.Domains, nil
}

// CheckSSLCertificateRequestDomains is equivalent to CheckSSLCertificateRequestDomainsContext with context.Background().
//
// Deprecated: use CheckSSLCertificateRequestDomainsContext.
func (api *API) CheckSSLCertificateRequestDomains(domains []string) (SSLCertificateRequestDomainChecks, error) {
	return api.CheckSSLCertificateRequestDomainsContext(context.Background(), domains)
}
