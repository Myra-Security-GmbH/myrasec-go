package myrasec

import (
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// TagInformation represents a key-value label attached to a resource.
type TagInformation struct {
	// ID is the unique identifier for the tag assignment.
	// This value is server-generated and required for update and delete operations.
	ID int `json:"id,omitempty" jsonschema:"The unique identifier for the tag. Server-generated; required for updates and deletes, but ignored during creation."`

	// Created indicates when the tag was added.
	// This is a server-managed, read-only value in ISO 8601 format.
	Created *types.DateTime `json:"created,omitempty" jsonschema:"The timestamp of creation (ISO 8601 format). Server-managed, read-only."`

	// Modified serves as a version identifier for optimistic locking.
	// It records the last update time in ISO 8601 format. This field is required
	// for update and delete operations to ensure data consistency.
	Modified *types.DateTime `json:"modified,omitempty" jsonschema:"The last update timestamp (ISO 8601 format). Required for updates and deletes to ensure data consistency (optimistic locking)."`

	// Key is the identifier or category name for the tag.
	// Example: "Environment" or "Project".
	Key string `json:"key" jsonschema:"The key portion of the tag pair (e.g., 'Environment')."`

	// Value is the specific content associated with the tag key.
	// Example: "Production" or "Marketing-Campaign".
	Value string `json:"value" jsonschema:"The value portion of the tag pair (e.g., 'Production')."`

	// Comment provides a descriptive note for this tag information key-value pair.
	Comment string `json:"comment,omitempty" jsonschema:"A descriptive comment or note regarding this specific tag information key-value pair."`
}

func getTagInformationMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"listTagInformation": {
			Name:   "listTagInformation",
			Action: "tags/%d/information",
			Method: http.MethodGet,
			Result: []TagInformation{},
		},
		"listTagInformationBySubDomainName": {
			Name:   "listTagInformationBySubDomainName",
			Action: "tags/%s/information",
			Method: http.MethodGet,
			Result: []TagInformation{},
		},
		"createTagInformation": {
			Name:   "createTagInformation",
			Action: "tags/%d/information",
			Method: http.MethodPost,
			Result: TagInformation{},
		},
		"updateTagInformation": {
			Name:   "updateTagInformation",
			Action: "tags/%d/information/%d",
			Method: http.MethodPut,
			Result: TagInformation{},
		},
		"deleteTagInformation": {
			Name:   "deleteTagInformation",
			Action: "tags/%d/information/%d",
			Method: http.MethodDelete,
			Result: TagInformation{},
		},
	}
}

// ListTagInformation returns a slice containing all tag information for the passed tag (ID)
func (api *API) ListTagInformation(tagId int, params map[string]string) ([]TagInformation, error) {
	if _, ok := api.methods["listTagInformation"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listTagInformation")
	}

	definition := api.methods["listTagInformation"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]TagInformation)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// ListTagInformationBySubDomainName returns a slice containing all tag information for the passed subDomainName
func (api *API) ListTagInformationBySubDomainName(subDomainName string, params map[string]string) ([]TagInformation, error) {
	if _, ok := api.methods["listTagInformationBySubDomainName"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listTagInformationBySubDomainName")
	}

	definition := api.methods["listTagInformationBySubDomainName"]
	definition.Action = fmt.Sprintf(definition.Action, subDomainName)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	res, ok := result.(*[]TagInformation)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return *res, nil
}

// CreateTagInformation creates a new tag information for the passed tag (ID) using the MYRA API
func (api *API) CreateTagInformation(information *TagInformation, tagId int) (*TagInformation, error) {
	if _, ok := api.methods["createTagInformation"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createTagInformation")
	}

	definition := api.methods["createTagInformation"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(definition, information)
	if err != nil {
		return nil, err
	}
	res, ok := result.(*TagInformation)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// UpdateTagInformation updates the passed tag information using the MYRA API
func (api *API) UpdateTagInformation(information *TagInformation, tagId int) (*TagInformation, error) {
	if _, ok := api.methods["updateTagInformation"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateTagInformation")
	}

	definition := api.methods["updateTagInformation"]
	definition.Action = fmt.Sprintf(definition.Action, tagId, information.ID)

	result, err := api.call(definition, information)
	if err != nil {
		return nil, err
	}
	res, ok := result.(*TagInformation)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// DeleteTagInformation deletes the passed tag information using the MYRA API
func (api *API) DeleteTagInformation(information *TagInformation, tagId int) (*TagInformation, error) {
	if _, ok := api.methods["deleteTagInformation"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteTagInformation")
	}

	definition := api.methods["deleteTagInformation"]
	definition.Action = fmt.Sprintf(definition.Action, tagId, information.ID)

	_, err := api.call(definition, information)
	if err != nil {
		return nil, err
	}
	return information, nil
}
