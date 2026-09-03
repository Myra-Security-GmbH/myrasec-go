package myrasec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// getDomainMethods returns Domain related API calls
func getDomainMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"getDomain": {
			Name:               "getDomain",
			Action:             "domains/%d",
			Method:             http.MethodGet,
			Result:             Domain{},
			ResponseDecodeFunc: decodeSingleElementResponse,
		},
		"listDomains": {
			Name:   "listDomains",
			Action: "domains",
			Method: http.MethodGet,
			Result: []Domain{},
		},
		"createDomain": {
			Name:   "createDomain",
			Action: "domains",
			Method: http.MethodPost,
			Result: Domain{},
		},
		"updateDomain": {
			Name:   "updateDomain",
			Action: "domains/%d",
			Method: http.MethodPut,
			Result: Domain{},
		},
		"deleteDomain": {
			Name:   "deleteDomain",
			Action: "domains/%d",
			Method: http.MethodDelete,
			Result: Domain{},
		},
	}
}

// Domain represents a website or zone managed by the Myra platform.
// It acts as the root object for all configuration settings, DNS records,
// and cache rules associated with a specific FQDN.
type Domain struct {
	// ID is the unique identifier for the domain.
	// This value is server-generated and required for update and delete operations.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the domain. Server-generated; required for updates and deletes, but ignored during creation."`

	// Created indicates when the domain was added to the system.
	// This is a server-managed, read-only value in ISO 8601 format.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified serves as a version identifier for optimistic locking.
	// It records the last update time in ISO 8601 format. This field is required
	// for update and delete operations to ensure data consistency.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Required for updates and deletes to ensure data consistency (optimistic locking)."`

	// Name is the Fully Qualified Domain Name (FQDN).
	// This value is immutable once created.
	Name string `json:"name" jsonschema:"The Fully Qualified Domain Name (FQDN) (e.g., 'example.com'). Immutable after creation. To correct a typo, the domain must be deleted and recreated."`

	// AutoUpdate controls the immediate deployment of configuration changes.
	// If false, changes are saved to the database but not propagated to the edge/WAF
	// until re-enabled. Useful for performing atomic batch updates.
	AutoUpdate bool `json:"autoUpdate" jsonschema:"Controls deployment of config changes. If false, changes are saved only to the database and NOT deployed to the edge network until re-enabled. Use this for atomic batch updates."`

	// AutoDNS triggers an automatic DNS record import during creation.
	// Note: This relies on external DNS queries and may result in an incomplete list.
	AutoDNS bool `json:"autoDns" jsonschema:"If true during creation, the system attempts to scan and import existing DNS records. Note: Results may be incomplete depending on the DNS provider."`

	// Paused indicates if the WAF/Protection is currently suspended for this domain.
	Paused bool `json:"paused" jsonschema:"Indicates if the Myra protection/WAF is currently suspended (paused) for this domain."`

	// PausedUntil specifies the scheduled date for automatic reactivation of protection.
	PausedUntil *types.DateTime `json:"pausedUntil,omitempty" jsonschema:"The timestamp (ISO 8601) when protection will be automatically reactivated. Only relevant if 'paused' is true."`

	// Reversed indicates if the domain is reversed.
	Reversed bool `json:"reversed" jsonschema:"Indicates whether the domain is reversed (boolean flag)."`
}

// GetDomainContext returns a single domain with/for the given identifier
func (api *API) GetDomainContext(ctx context.Context, id int) (*Domain, error) {
	if _, ok := api.methods["getDomain"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getDomain")
	}

	definition := api.methods["getDomain"]
	definition.Action = fmt.Sprintf(definition.Action, id)

	result, err := api.call(ctx, definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	res, ok := result.(*Domain)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// GetDomain is equivalent to GetDomainContext with context.Background().
//
// Deprecated: use GetDomainContext.
func (api *API) GetDomain(id int) (*Domain, error) {
	return api.GetDomainContext(context.Background(), id)
}

// ListDomainsContext returns a slice containing all visible domains
func (api *API) ListDomainsContext(ctx context.Context, params map[string]string) ([]Domain, error) {
	if _, ok := api.methods["listDomains"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listDomains")
	}

	definition := api.methods["listDomains"]

	result, err := api.call(ctx, definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]Domain)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// ListDomains is equivalent to ListDomainsContext with context.Background().
//
// Deprecated: use ListDomainsContext.
func (api *API) ListDomains(params map[string]string) ([]Domain, error) {
	return api.ListDomainsContext(context.Background(), params)
}

// CreateDomainContext creates a new domain using the MYRA API
func (api *API) CreateDomainContext(ctx context.Context, domain *Domain) (*Domain, error) {
	if _, ok := api.methods["createDomain"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createDomain")
	}

	definition := api.methods["createDomain"]

	result, err := api.call(ctx, definition, domain)
	if err != nil {
		return nil, err
	}
	res, ok := result.(*Domain)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// CreateDomain is equivalent to CreateDomainContext with context.Background().
//
// Deprecated: use CreateDomainContext.
func (api *API) CreateDomain(domain *Domain) (*Domain, error) {
	return api.CreateDomainContext(context.Background(), domain)
}

// UpdateDomainContext updates the passed domain using the MYRA API
func (api *API) UpdateDomainContext(ctx context.Context, domain *Domain) (*Domain, error) {
	if _, ok := api.methods["updateDomain"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateDomain")
	}

	definition := api.methods["updateDomain"]
	definition.Action = fmt.Sprintf(definition.Action, domain.ID)

	result, err := api.call(ctx, definition, domain)
	if err != nil {
		return nil, err
	}
	res, ok := result.(*Domain)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// UpdateDomain is equivalent to UpdateDomainContext with context.Background().
//
// Deprecated: use UpdateDomainContext.
func (api *API) UpdateDomain(domain *Domain) (*Domain, error) {
	return api.UpdateDomainContext(context.Background(), domain)
}

// DeleteDomainContext deletes the passed domain using the MYRA API
func (api *API) DeleteDomainContext(ctx context.Context, domain *Domain) (*Domain, error) {
	if _, ok := api.methods["deleteDomain"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteDomain")
	}

	definition := api.methods["deleteDomain"]
	definition.Action = fmt.Sprintf(definition.Action, domain.ID)

	_, err := api.call(ctx, definition, domain)
	if err != nil {
		return nil, err
	}
	return domain, nil
}

// DeleteDomain is equivalent to DeleteDomainContext with context.Background().
//
// Deprecated: use DeleteDomainContext.
func (api *API) DeleteDomain(domain *Domain) (*Domain, error) {
	return api.DeleteDomainContext(context.Background(), domain)
}
