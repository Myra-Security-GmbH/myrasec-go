package myrasec

import (
	"fmt"
)

//
// ListTagSettings returns a Setting struct containing the settings for the passed tag
//
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

//
// UpdateTagSettings updates the passed settings using the MYRA API
//
func (api *API) UpdateTagSettings(settings *Settings, tagId int) (*Settings, error) {
	if _, ok := methods["updateTagSettings"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createTagSettings")
	}

	definition := methods["updateTagSettings"]
	definition.Action = fmt.Sprintf(definition.Action, tagId)

	result, err := api.call(definition, settings)
	if err != nil {
		return nil, err
	}
	return result.(*Settings), nil

}
