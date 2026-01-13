package myrasec

import (
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// getTagMethods returns Tag related API calls
func getTagMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"getTag": {
			Name:               "getTag",
			Action:             "tags/%d",
			Method:             http.MethodGet,
			Result:             Tag{},
			ResponseDecodeFunc: decodeSingleElementResponse,
		},
		"listTags": {
			Name:   "listTags",
			Action: "tags",
			Method: http.MethodGet,
			Result: []Tag{},
		},
		"createTag": {
			Name:   "createTag",
			Action: "tags",
			Method: http.MethodPost,
			Result: Tag{},
		},
		"updateTag": {
			Name:   "updateTag",
			Action: "tags/%d",
			Method: http.MethodPut,
			Result: Tag{},
		},
		"deleteTag": {
			Name:   "deleteTag",
			Action: "tags/%d",
			Method: http.MethodDelete,
			Result: Tag{},
		},
		"cloneTag": {
			Name:   "cloneTag",
			Action: "tags/%d",
			Method: http.MethodPost,
			Result: Tag{},
		},
	}
}

// Tag represents a logical grouping container for configuration settings.
// Tags can be applied to domains or subdomains to enforce shared rules (e.g., WAF, Cache).
type Tag struct {
	// ID is the unique identifier for the tag.
	// This value is server-generated and required for update and delete operations.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the tag. Server-generated; required for updates and deletes, but ignored during creation."`

	// Created indicates when the tag was created.
	// This is a server-managed, read-only value in ISO 8601 format.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified serves as a version identifier for optimistic locking.
	// It records the last update time in ISO 8601 format. This field is required
	// for update and delete operations to ensure data consistency.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Required for updates and deletes to ensure data consistency (optimistic locking)."`

	// Name is the display label for the tag.
	Name string `json:"name" jsonschema:"The unique name of the tag used for identification."`

	// Type defines the functional category of the tag.
	// Valid values are: CONFIG, WAF, CACHE, RATE_LIMIT, INFORMATION.
	Type string `json:"type" jsonschema:"The category/type of the tag. Valid values: 'CONFIG', 'WAF', 'CACHE', 'RATE_LIMIT', 'INFORMATION'."`

	// Assignments lists the domains or subdomains linked to this tag.
	Assignments []TagAssignment `json:"assignments" jsonschema:"List of resources (domains/subdomains) assigned to this tag."`

	// Sort defines the processing order priority.
	// Specifically relevant for WAF tags to determine rule execution order.
	Sort int `json:"sort,omitempty" jsonschema:"Priority/Sorting order. Crucial for 'WAF' tags to determine the execution order of rules."`

	// Global indicates if the tag is a system-wide predefined tag.
	// Global tags cannot be renamed or modified; only assignments can be changed.
	Global bool `json:"global,omitempty" jsonschema:"If true, this is a system-managed global tag. The tag definition (Name, Type) is read-only/immutable. You can only modify the 'assignments' list."`
}

// TagAssignment represents the link between a Tag and a specific resource (Domain or Subdomain).
type TagAssignment struct {
	// ID is the unique identifier for the assignment.
	// This value is server-generated and required for update and delete operations.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the assignment. Server-generated; required for updates and deletes, but ignored during creation."`

	// Created indicates when the assignment was created.
	// This is a server-managed, read-only value in ISO 8601 format.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified serves as a version identifier for optimistic locking.
	// It records the last update time in ISO 8601 format. This field is required
	// for update and delete operations to ensure data consistency.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Required for updates and deletes to ensure data consistency (optimistic locking)."`

	// Type defines the scope of the assignment.
	// Valid values are: DOMAIN, SUBDOMAIN.
	Type string `json:"type" jsonschema:"The scope of the assignment. Valid values: 'DOMAIN', 'SUBDOMAIN'."`

	// Title is the identifier of the assigned resource (usually the domain name).
	Title string `json:"title" jsonschema:"The identifier of the assigned resource (e.g., the Domain Name)."`

	// SubDomainName specifies the target subdomain.
	// Only required if Type is set to SUBDOMAIN.
	SubDomainName string `json:"subDomainName" jsonschema:"The specific subdomain FQDN. Required if 'type' is set to 'SUBDOMAIN'."`
}

// GetTag returns a single tag for the given identifier
func (api *API) GetTag(id int) (*Tag, error) {
	if _, ok := methods["getTag"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getTag")
	}

	definition := methods["getTag"]
	definition.Action = fmt.Sprintf(definition.Action, id)

	result, err := api.call(definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	return result.(*Tag), nil
}

// ListTags returns a slice containing all visible tags
func (api *API) ListTags(params map[string]string) ([]Tag, error) {
	if _, ok := methods["listTags"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listTags")
	}

	definition := methods["listTags"]

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	return *result.(*[]Tag), nil
}

// CreateTag creates a new tag using the MYRA API
func (api *API) CreateTag(tag *Tag) (*Tag, error) {
	if _, ok := methods["createTag"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createTag")
	}

	definition := methods["createTag"]

	result, err := api.call(definition, tag)
	if err != nil {
		return nil, err
	}

	return result.(*Tag), nil
}

// UpdateTag updates the passed tag using the MYRA API
func (api *API) UpdateTag(tag *Tag) (*Tag, error) {
	if _, ok := methods["updateTag"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateTag")
	}

	definition := methods["updateTag"]
	definition.Action = fmt.Sprintf(definition.Action, tag.ID)

	result, err := api.call(definition, tag)
	if err != nil {
		return nil, err
	}

	return result.(*Tag), nil
}

// DeleteTag deletes the passed tag using the MYRA API
func (api *API) DeleteTag(tag *Tag) (*Tag, error) {
	if _, ok := methods["deleteTag"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteTag")
	}

	definition := methods["deleteTag"]
	definition.Action = fmt.Sprintf(definition.Action, tag.ID)

	_, err := api.call(definition, tag)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

// CloneTag clones the passed tag using the MYRA API
func (api *API) CloneTag(tag *Tag) (*Tag, error) {
	if _, ok := methods["cloneTag"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "cloneTag")
	}

	definition := methods["cloneTag"]
	definition.Action = fmt.Sprintf(definition.Action, tag.ID)

	result, err := api.call(definition, tag)
	if err != nil {
		return nil, err
	}

	return result.(*Tag), nil
}
