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

// Tag ...
type Tag struct {
	ID          int             `json:"id,omitempty" jsonschema:"ID is an unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete a Tag it is necessary to add this attribute to your object."`
	Created     *types.DateTime `json:"created,omitempty" jsonschema:"Created is a date type attribute with an ISO 8601 format. Created will be created by the server after creating a new Tag object. This value is informational so it is not necessary to add this attribute to any API call."`
	Modified    *types.DateTime `json:"modified,omitempty" jsonschema:"Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates and deletes. This value is always a date type with an ISO 8601 format."`
	Name        string          `json:"name" jsonschema:"Identifies the tag by its name."`
	Type        string          `json:"type" jsonschema:"Defines the type of the tag and must be one of CONFIG, WAF, CACHE, RATE_LIMIT, INFORMATION ."`
	Assignments []TagAssignment `json:"assignments"`
	Sort        int             `json:"sort,omitempty" jsonschema:"Defines the order in which WAF tags are processed."`
	Global      bool            `json:"global,omitempty" jsonschema:"Identify global tags. It is not possible to edit a global tags. It is only possible to assign (sub)domains to a global tag."`
}

// TagAssignment ...
type TagAssignment struct {
	ID            int             `json:"id,omitempty" jsonschema:"ID is an unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete a TagAssignment it is necessary to add this attribute to your object."`
	Created       *types.DateTime `json:"created,omitempty" jsonschema:"Created is a date type attribute with an ISO 8601 format. Created will be created by the server after creating a new TagAssignment object. This value is informational so it is not necessary to add this attribute to any API call."`
	Modified      *types.DateTime `json:"modified,omitempty" jsonschema:"Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates and deletes. This value is always a date type with an ISO 8601 format."`
	Type          string          `json:"type" jsonschema:"Defines the type of the tag assignment and must be one of DOMAIN, SUBDOMAIN ."`
	Title         string          `json:"title" jsonschema:"Identifies the tag assignment by its domain name."`
	SubDomainName string          `json:"subDomainName" jsonschema:"Only set on SUBDOMAIN tag assignments."`
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
