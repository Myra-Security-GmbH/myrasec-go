package myrasec

import (
	"fmt"
	"net/http"
)

func getTagInformationMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"getTagInformation": {
			Name:               "getTagInformation",
			Action:             "tag/information/%d",
			Method:             http.MethodGet,
			Result:             map[string]interface{}{},
			ResponseDecodeFunc: decodeSingleElementResponse,
		},
		"listTagInformationBySubDomainName": {
			Name:   "listTagInformationBySubDomainName",
			Action: "tag/information/%s",
			Method: http.MethodGet,
			Result: map[string]interface{}{},
		},
		"updateTagInformation": {
			Name:   "updateTagInformation",
			Action: "tag/information/%d",
			Method: http.MethodPost,
			Result: map[string]interface{}{},
		},
	}
}

func (api *API) GetTagInformation(tagId int) (map[string]interface{}, error) {
	if _, ok := methods["getTagInformation"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getTagInformation")
	}

	definition := methods["getTagInformation"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	return *result.(*map[string]interface{}), nil
}

func (api *API) ListTagInformationsBySubDomainName(subDomainName string) (map[string]interface{}, error) {
	if _, ok := methods["listTagInformationsBySubDomainName"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listTagInformationBySubDomainName")
	}

	definition := methods["listTagInformationBySubDomainName"]
	definition.Action = fmt.Sprintf(definition.Action, subDomainName)

	result, err := api.call(definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	return *result.(*map[string]interface{}), nil
}

func (api *API) UpdateTagInformation(tagId int, informations map[string]interface{}) (map[string]interface{}, error) {
	if _, ok := methods["updateTagInformation"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateTagInformation")
	}

	definition := methods["updateTagInformation"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(definition, informations)
	if err != nil {
		return nil, err
	}

	return *result.(*map[string]interface{}), nil
}
