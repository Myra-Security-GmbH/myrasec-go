package myrasec

import (
	"fmt"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
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
	MatchingType  string          `json:"matchingType"`
	Comment       string          `json:"comment,omitempty"`
	Sort          int             `json:"sort,omitempty"`
	Enabled       bool            `json:"enabled"`
	ExpertMode    bool            `json:"expertMode,omitempty"`
}

//
// GetRedirect returns a single redirect with/for the given identifier
//
func (api *API) GetRedirect(domainId int, subDomainName string, id int) (*Redirect, error) {
	if _, ok := methods["getRedirect"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getRedirect")
	}

	definition := methods["getRedirect"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, id)

	result, err := api.call(definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	return result.(*Redirect), nil
}

//
// ListRedirects returns a slice containing all visible redirects for a subdomain
//
func (api *API) ListRedirects(domainId int, subDomainName string, params map[string]string) ([]Redirect, error) {
	if _, ok := methods["listRedirects"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listRedirects")
	}

	definition := methods["listRedirects"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}
	var records []Redirect
	records = append(records, *result.(*[]Redirect)...)

	return records, nil
}

//
// CreateRedirect creates a new redirect for the passed subdomain (name) using the MYRA API
//
func (api *API) CreateRedirect(redirect *Redirect, domainId int, subDomainName string) (*Redirect, error) {
	if _, ok := methods["createRedirect"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createRedirect")
	}

	definition := methods["createRedirect"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName)

	result, err := api.call(definition, redirect)
	if err != nil {
		return nil, err
	}
	return result.(*Redirect), nil
}

//
// UpdateRedirect updates the passed redirect using the MYRA API
//
func (api *API) UpdateRedirect(redirect *Redirect, domainId int, subDomainName string) (*Redirect, error) {
	if _, ok := methods["updateRedirect"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateRedirect")
	}

	definition := methods["updateRedirect"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, redirect.ID)

	result, err := api.call(definition, redirect)
	if err != nil {
		return nil, err
	}
	return result.(*Redirect), nil
}

//
// DeleteRedirect deletes the passed redirect using the MYRA API
//
func (api *API) DeleteRedirect(redirect *Redirect, domainId int, subDomainName string) (*Redirect, error) {
	if _, ok := methods["deleteRedirect"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteRedirect")
	}

	definition := methods["deleteRedirect"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, subDomainName, redirect.ID)

	_, err := api.call(definition, redirect)
	if err != nil {
		return nil, err
	}
	return redirect, nil
}
