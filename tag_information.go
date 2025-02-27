package myrasec

import (
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

type TagInformation struct {
	ID       int             `json:"id,omitempty"`
	Created  *types.DateTime `json:"created,omitempty"`
	Modified *types.DateTime `json:"modified,omitempty"`
	Key      string          `json:"key"`
	Value    string          `json:"value"`
	Comment  string          `json:"comment,omitempty"`
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
	if _, ok := methods["listTagInformation"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listTagInformation")
	}

	definition := methods["listTagInformation"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	return *result.(*[]TagInformation), nil
}

// ListTagInformationBySubDomainName returns a slice containing all tag information for the passed subDomainName
func (api *API) ListTagInformationBySubDomainName(subDomainName string, params map[string]string) ([]TagInformation, error) {
	if _, ok := methods["listTagInformationBySubDomainName"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listTagInformationBySubDomainName")
	}

	definition := methods["listTagInformationBySubDomainName"]
	definition.Action = fmt.Sprintf(definition.Action, subDomainName)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	return *result.(*[]TagInformation), nil

}

// CreateTagInformation creates a new tag information for the passed tag (ID) using the MYRA API
func (api *API) CreateTagInformation(information *TagInformation, tagId int) (*TagInformation, error) {
	if _, ok := methods["createTagInformation"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createTagInformation")
	}

	definition := methods["createTagInformation"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(definition, information)
	if err != nil {
		return nil, err
	}
	return result.(*TagInformation), nil
}

// UpdateTagInformation updates the passed tag information using the MYRA API
func (api *API) UpdateTagInformation(information *TagInformation, tagId int) (*TagInformation, error) {
	if _, ok := methods["updateTagInformation"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateTagInformation")
	}

	definition := methods["updateTagInformation"]
	definition.Action = fmt.Sprintf(definition.Action, tagId, information.ID)

	result, err := api.call(definition, information)
	if err != nil {
		return nil, err
	}
	return result.(*TagInformation), nil
}

// DeleteTagInformation deletes the passed tag information using the MYRA API
func (api *API) DeleteTagInformation(information *TagInformation, tagId int) (*TagInformation, error) {
	if _, ok := methods["deleteTagInformation"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteTagInformation")
	}

	definition := methods["deleteTagInformation"]
	definition.Action = fmt.Sprintf(definition.Action, tagId, information.ID)

	_, err := api.call(definition, information)
	if err != nil {
		return nil, err
	}
	return information, nil
}
