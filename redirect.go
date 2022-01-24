package myrasec

import (
	"fmt"
	"strconv"

	"github.com/Myra-Security-GmbH/myrasec-go/pkg/types"
)

//
// Redirect ...
//
type Redirect struct {
	ID            int             `json:"id,omitempty"`
	Created       *types.DateTime `json:"created,omitempty"`
	Modified      *types.DateTime `json:"modified,omitempty"`
	Type          string          `json:"type"`
	SubDomainName string          `json:"subDomainName"`
	Source        string          `json:"source"`
	Destination   string          `json:"destination"`
	Sort          int             `json:"sort,omitempty"`
	MatchingType  string          `json:"matchingType"`
	Enabled       bool            `json:"enabled,omitempty"`
}

//
// ListRedirects returns a slice containing all visible redirects for a subdomain
//
func (api *API) ListRedirects(subDomainName string, params map[string]string) ([]Redirect, error) {
	if _, ok := methods["listRedirects"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "listRedirects")
	}

	page := 1
	var err error
	if pageParam, ok := params[ParamPage]; ok {
		delete(params, ParamPage)
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			page = 1
		}
	}

	definition := methods["listRedirects"]
	definition.Action = fmt.Sprintf(definition.Action, subDomainName, page)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}
	var records []Redirect
	for _, v := range *result.(*[]Redirect) {
		records = append(records, v)
	}

	return records, nil
}

//
// CreateRedirect creates a new redirect for the passed subdomain (name) using the MYRA API
//
func (api *API) CreateRedirect(redirect *Redirect, subDomainName string) (*Redirect, error) {
	if _, ok := methods["createRedirect"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "createRedirect")
	}

	definition := methods["createRedirect"]
	definition.Action = fmt.Sprintf(definition.Action, subDomainName)

	result, err := api.call(definition, redirect)
	if err != nil {
		return nil, err
	}
	return result.(*Redirect), nil
}

//
// UpdateRedirect updates the passed redirect using the MYRA API
//
func (api *API) UpdateRedirect(redirect *Redirect, subDomainName string) (*Redirect, error) {
	if _, ok := methods["updateRedirect"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "updateRedirect")
	}

	definition := methods["updateRedirect"]
	definition.Action = fmt.Sprintf(definition.Action, subDomainName)

	result, err := api.call(definition, redirect)
	if err != nil {
		return nil, err
	}
	return result.(*Redirect), nil
}

//
// DeleteRedirect deletes the passed redirect using the MYRA API
//
func (api *API) DeleteRedirect(redirect *Redirect, subDomainName string) (*Redirect, error) {
	if _, ok := methods["deleteRedirect"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "deleteRedirect")
	}

	definition := methods["deleteRedirect"]
	definition.Action = fmt.Sprintf(definition.Action, subDomainName)

	result, err := api.call(definition, redirect)
	if err != nil {
		return nil, err
	}
	return result.(*Redirect), nil
}
