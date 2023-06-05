package myrasec

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// getTagSettingsMethods returns Tag settings related API calls
func getTagSettingsMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"listTagSettings": {
			Name:               "listTagSettings",
			Action:             "tag/%d/settings",
			Method:             http.MethodGet,
			Result:             Settings{},
			ResponseDecodeFunc: decodeTagSettingsResponse,
		},
		"listTagSettingsMap": {
			Name:               "listTagSettingsMap",
			Action:             "tag/%d/settings",
			Method:             http.MethodGet,
			Result:             map[string]interface{}{},
			ResponseDecodeFunc: decodeTagSettingsMapResponse,
		},
		"updateTagSettings": {
			Name:   "updateTagSettings",
			Action: "tag/%d/settings",
			Method: http.MethodPut,
			Result: Settings{},
		},
		"updateTagSettingsPartial": {
			Name:   "updateTagSettingsPartial",
			Action: "tag/%d/settings",
			Method: http.MethodPut,
			Result: map[string]interface{}{},
		},
	}
}

// tagSettingsResponse ...
type tagSettingsResponse struct {
	Settings Settings `json:"settings"`
}

// ListTagSettings returns a Setting struct containing the settings for the passed tag
func (api *API) ListTagSettings(tagId int) (*Settings, error) {
	if _, ok := methods["listTagSettings"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listTagSettings")
	}

	definition := methods["listTagSettings"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	return result.(*Settings), nil
}

func (api *API) ListTagSettingsMap(tagId int) (interface{}, error) {
	if _, ok := methods["listTagSettingsMap"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listTagSettingsMap")
	}

	definition := methods["listTagSettingsMap"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(definition, map[string]string{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateTagSettings updates the passed settings using the MYRA API
func (api *API) UpdateTagSettings(settings *Settings, tagId int) (*Settings, error) {
	if _, ok := methods["updateTagSettings"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateTagSettings")
	}

	definition := methods["updateTagSettings"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(definition, settings)
	if err != nil {
		return nil, err
	}
	return result.(*Settings), nil

}

// UpdateTagSettings updates the passed settings using the MYRA API
func (api *API) UpdateTagSettingsPartial(settings map[string]interface{}, tagId int) (interface{}, error) {
	if _, ok := methods["updateTagSettingsPartial"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateTagSettingsPartial")
	}

	definition := methods["updateTagSettingsPartial"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(definition, settings)
	if err != nil {
		return nil, err
	}
	return result, nil

}

// decodeTagSettingsResponse - custom decode function for tag settings response. Used in the ListTagSettings action.
func decodeTagSettingsResponse(resp *http.Response, definition APIMethod) (interface{}, error) {
	var res tagSettingsResponse
	err := json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res.Settings, nil
}

// decodeSettingsResponseFull - custom decode function for full settings response. Used in the ListSettingsFull action.
func decodeTagSettingsMapResponse(resp *http.Response, definition APIMethod) (interface{}, error) {
	var res map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
