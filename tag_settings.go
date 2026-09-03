package myrasec

import (
	"context"
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
			Result:             map[string]any{},
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
			Result: map[string]any{},
		},
	}
}

// tagSettingsResponse wraps the Settings struct for tag settings API responses.
type tagSettingsResponse struct {
	Settings Settings `json:"settings"`
}

// ListTagSettingsContext returns a Setting struct containing the settings for the passed tag
func (api *API) ListTagSettingsContext(ctx context.Context, tagId int) (*Settings, error) {
	if _, ok := api.methods["listTagSettings"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listTagSettings")
	}

	definition := api.methods["listTagSettings"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(ctx, definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	res, ok := result.(*Settings)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// ListTagSettings is equivalent to ListTagSettingsContext with context.Background().
//
// Deprecated: use ListTagSettingsContext.
func (api *API) ListTagSettings(tagId int) (*Settings, error) {
	return api.ListTagSettingsContext(context.Background(), tagId)
}

func (api *API) ListTagSettingsMapContext(ctx context.Context, tagId int) (any, error) {
	if _, ok := api.methods["listTagSettingsMap"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listTagSettingsMap")
	}

	definition := api.methods["listTagSettingsMap"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(ctx, definition, map[string]string{})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ListTagSettingsMap is equivalent to ListTagSettingsMapContext with context.Background().
//
// Deprecated: use ListTagSettingsMapContext.
func (api *API) ListTagSettingsMap(tagId int) (any, error) {
	return api.ListTagSettingsMapContext(context.Background(), tagId)
}

// UpdateTagSettingsContext updates the passed settings using the MYRA API
func (api *API) UpdateTagSettingsContext(ctx context.Context, settings *Settings, tagId int) (*Settings, error) {
	if _, ok := api.methods["updateTagSettings"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateTagSettings")
	}

	definition := api.methods["updateTagSettings"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(ctx, definition, settings)
	if err != nil {
		return nil, err
	}
	res, ok := result.(*Settings)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
}

// UpdateTagSettings is equivalent to UpdateTagSettingsContext with context.Background().
//
// Deprecated: use UpdateTagSettingsContext.
func (api *API) UpdateTagSettings(settings *Settings, tagId int) (*Settings, error) {
	return api.UpdateTagSettingsContext(context.Background(), settings, tagId)
}

// UpdateTagSettings updates the passed settings using the MYRA API
func (api *API) UpdateTagSettingsPartialContext(ctx context.Context, settings map[string]any, tagId int) (any, error) {
	if _, ok := api.methods["updateTagSettingsPartial"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateTagSettingsPartial")
	}

	definition := api.methods["updateTagSettingsPartial"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(ctx, definition, settings)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateTagSettingsPartial is equivalent to UpdateTagSettingsPartialContext with context.Background().
//
// Deprecated: use UpdateTagSettingsPartialContext.
func (api *API) UpdateTagSettingsPartial(settings map[string]any, tagId int) (any, error) {
	return api.UpdateTagSettingsPartialContext(context.Background(), settings, tagId)
}

// decodeTagSettingsResponse - custom decode function for tag settings response. Used in the ListTagSettings action.
func decodeTagSettingsResponse(resp *http.Response, definition APIMethod) (any, error) {
	var res tagSettingsResponse
	err := json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res.Settings, nil
}

// decodeSettingsResponseFull - custom decode function for full settings response. Used in the ListSettingsFull action.
func decodeTagSettingsMapResponse(resp *http.Response, definition APIMethod) (any, error) {
	var res map[string]any
	err := json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
