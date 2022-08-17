package myrasec

import (
	"fmt"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

//
// Tag
//
type Tag struct {
	ID           int             `json:"id,omitempty"`
	Created      *types.DateTime `json:"created,omitempty"`
	Modified     *types.DateTime `json:"modified,omitempty"`
	Name         string          `json:"name"`
	Type         string          `json:"type"`
	Organization int             `json:"organization"`
	Assignments  []TagAssignment `json:"assignments"`
}

type TagAssignment struct {
	ID            int             `json:"id,omitempty"`
	Created       *types.DateTime `json:"created,omitempty"`
	Modified      *types.DateTime `json:"modified,omitempty"`
	Type          string          `json:"type"`
	Title         string          `json:"title"`
	SubDomainName string          `json:"subDomainName"`
}

//
// GetTag returns a single tag for the given identifier
//
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

//
// ListTags returns a slice containing all visible tags
//
func (api *API) ListTags(params map[string]string) ([]Tag, error) {
	if _, ok := methods["listTags"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listTags")
	}

	definition := methods["listTags"]

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	var tags []Tag
	tags = append(tags, *result.(*[]Tag)...)

	return tags, nil
}

//
// CreateTag creates a new tag using the MYRA API
//
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

//
// UpdateTag updates the passed tag using the MYRA API
//
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

//
// DeleteTag deletes the passed tag using the MYRA API
//
func (api *API) DeleteTag(tag *Tag) (*Tag, error) {
	if _, ok := methods["deleteTag"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteTag")
	}

	definition := methods["deleteTag"]
	definition.Action = fmt.Sprintf(definition.Action, tag.ID)

	result, err := api.call(definition, tag)
	if err != nil {
		return nil, err
	}

	return result.(*Tag), nil

}
