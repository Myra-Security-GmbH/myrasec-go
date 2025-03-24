package myrasec

import (
	"fmt"
	"net/http"
)

func getZoneConfigMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"getZoneConfigRaw": {
			Name:   "getZoneConfigRaw",
			Action: "domains/%d/bind-raw/%s",
			Method: http.MethodGet,
		},
		"getZoneConfigJson": {
			Name:   "getZoneConfigJson",
			Action: "domains/%d/bind/%s",
			Method: http.MethodGet,
		},
	}
}

func (api *API) GetZoneConfigRaw(domainId int, params map[string]string) ([]byte, error) {
	if _, ok := methods["getZoneConfigRaw"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getZoneConfigRaw")
	}

	ipTarget, ok := params["ipTarget"]
	if !ok {
		ipTarget = "myra"
	}
	definition := methods["getZoneConfigRaw"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, ipTarget)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	return result.([]byte), nil
}

func (api *API) GetZoneConfigJson(domainId int, params map[string]string) ([]byte, error) {
	if _, ok := methods["getZoneConfigJson"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getZoneConfigJson")
	}

	ipTarget, ok := params["ipTarget"]
	if !ok {
		ipTarget = "myra"
	}
	definition := methods["getZoneConfigJson"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, ipTarget)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	return result.([]byte), nil

}
