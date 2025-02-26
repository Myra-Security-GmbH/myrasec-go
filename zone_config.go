package myrasec

import (
	"fmt"
	"io"
	"net/http"
)

func getZoneConfigMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"getZoneConfigRaw": {
			Name:               "getZoneConfigRaw",
			Action:             "domains/%d/bind-raw/%s",
			Method:             http.MethodGet,
			ResponseDecodeFunc: decodeStringValue,
		},
		"getZoneConfigJson": {
			Name:               "getZoneConfigJson",
			Action:             "domains/%d/bind/%s",
			Method:             http.MethodGet,
			ResponseDecodeFunc: decodeStringValue,
		},
	}
}

func (api *API) GetZoneConfigRaw(domainId int, params map[string]string) (string, error) {
	if _, ok := methods["getZoneConfigRaw"]; !ok {
		return "", fmt.Errorf("passed action [%s] is not supported", "getZoneConfigRaw")
	}

	ipTarget, ok := params["ipTarget"]
	if !ok {
		ipTarget = "myra"
	}
	definition := methods["getZoneConfigRaw"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, ipTarget)

	result, err := api.call(definition, params)
	if err != nil {
		return "", err
	}

	return result.(string), nil
}

func (api *API) GetZoneConfigJson(domainId int, params map[string]string) (string, error) {
	if _, ok := methods["getZoneConfigJson"]; !ok {
		return "", fmt.Errorf("passed action [%s] is not supported", "getZoneConfigJson")
	}

	ipTarget, ok := params["ipTarget"]
	if !ok {
		ipTarget = "myra"
	}
	definition := methods["getZoneConfigJson"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, ipTarget)

	result, err := api.call(definition, params)
	if err != nil {
		return "", err
	}

	return result.(string), nil

}

func decodeStringValue(resp *http.Response, definition APIMethod) (interface{}, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(body[:])
	return bodyString, nil
}
