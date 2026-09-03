package myrasec

import (
	"context"
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

func (api *API) GetZoneConfigRawContext(ctx context.Context, domainId int, params map[string]string) (*string, error) {
	if _, ok := api.methods["getZoneConfigRaw"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getZoneConfigRaw")
	}

	ipTarget, ok := params["ipTarget"]
	if !ok {
		ipTarget = "myra"
	}
	definition := api.methods["getZoneConfigRaw"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, ipTarget)

	result, err := api.call(ctx, definition, params)
	if err != nil {
		return nil, err
	}

	res := result.(string)
	return &res, nil
}

// GetZoneConfigRaw is equivalent to GetZoneConfigRawContext with context.Background().
//
// Deprecated: use GetZoneConfigRawContext.
func (api *API) GetZoneConfigRaw(domainId int, params map[string]string) (*string, error) {
	return api.GetZoneConfigRawContext(context.Background(), domainId, params)
}

func (api *API) GetZoneConfigJsonContext(ctx context.Context, domainId int, params map[string]string) (*string, error) {
	if _, ok := api.methods["getZoneConfigJson"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getZoneConfigJson")
	}

	ipTarget, ok := params["ipTarget"]
	if !ok {
		ipTarget = "myra"
	}
	definition := api.methods["getZoneConfigJson"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, ipTarget)

	result, err := api.call(ctx, definition, params)
	if err != nil {
		return nil, err
	}

	ret := result.(string)
	return &ret, nil
}

// GetZoneConfigJson is equivalent to GetZoneConfigJsonContext with context.Background().
//
// Deprecated: use GetZoneConfigJsonContext.
func (api *API) GetZoneConfigJson(domainId int, params map[string]string) (*string, error) {
	return api.GetZoneConfigJsonContext(context.Background(), domainId, params)
}

func decodeStringValue(resp *http.Response, definition APIMethod) (any, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(body[:])
	return bodyString, nil
}
