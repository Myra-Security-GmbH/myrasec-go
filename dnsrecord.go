package myrasec

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/Myra-Security-GmbH/myrasec-go/pkg/types"
)

//
// DNSRecord ...
//
type DNSRecord struct {
	ID               int              `json:"id,omitempty"`
	Created          *types.DateTime  `json:"created,omitempty"`
	Modified         *types.DateTime  `json:"modified,omitempty"`
	Name             string           `json:"name"`
	Value            string           `json:"value"`
	RecordType       string           `json:"recordType"`
	AlternativeCNAME string           `json:"alternativeCname,omitempty"`
	Comment          string           `json:"comment,omitempty"`
	Active           bool             `json:"active"`
	Enabled          bool             `json:"enabled"`
	TTL              int              `json:"ttl"`
	Priority         int              `json:"priority,omitempty"`
	Port             int              `json:"port,omitempty"`
	UpstreamOptions  *UpstreamOptions `json:"upstreamOptions,omitempty"`
}

//
// UpstreamOptions ...
//
type UpstreamOptions struct {
	ID          int             `json:"id,omitempty"`
	Created     *types.DateTime `json:"created,omitempty"`
	Modified    *types.DateTime `json:"modified,omitempty"`
	Backup      bool            `json:"backup"`
	Down        bool            `json:"down"`
	FailTimeout int             `json:"failTimeout"`
	MaxFails    int             `json:"maxFails"`
	Weight      int             `json:"weight"`
}

//
// ListDNSRecords returns a slice containing all visible DNS records for a domain
//
func (api *API) ListDNSRecords(domainName string, params map[string]string) (Output, error) {
	var output Output

	if _, ok := methods["listDNSRecords"]; !ok {
		return output, fmt.Errorf("passed action [%s] is not supported", "listDNSRecords")
	}

	pageNumber, _ := strconv.ParseInt(params["pageNumber"], 10, 64)
	if pageNumber == 0 {
		pageNumber = 1
	}

	pageSize, _ := strconv.ParseInt(params["pageSize"], 10, 64)
	if pageSize == 0 {
		pageSize = 50 // default pageSize on API side
	}

	definition := methods["listDNSRecords"]
	definition.Action = fmt.Sprintf(definition.Action, domainName, pageNumber)

	response, err := api.call(definition, params)

	//@TODO lets replace it by some correct solution
	var ok bool
	output, ok = response.(Output)
	if ok == false {
		fmt.Errorf("casting response failed, interface not compatible")
	}

	if err != nil {
		return output, err
	}

	return output, nil
}

func decodeListDNSRecordsResponse(resp *http.Response, definition APIMethod) (interface{}, error) {
	var decodedResponse Response
	err := json.NewDecoder(resp.Body).Decode(&decodedResponse)
	if err != nil {
		return nil, err
	}

	if decodedResponse.Error {
		var errorMsg string
		for _, e := range decodedResponse.ViolationList {
			errorMsg += fmt.Sprintf("%s: %s\n", e.Path, e.Message)
		}
		return nil, errors.New(errorMsg)
	}

	var result []interface{}

	if decodedResponse.List != nil {
		result = decodedResponse.List
	}

	tmp, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(bytes.NewReader(tmp))

	retValue := reflect.New(reflect.TypeOf([]DNSRecord{}))
	res := retValue.Interface()
	decoder.Decode(res)

	output := Output{
		Count:    decodedResponse.Count,
		Page:     decodedResponse.Page,
		PageSize: decodedResponse.PageSize,
	}

	for _, v := range *res.(*[]DNSRecord) {
		output.Elements = append(output.Elements, v)
	}

	return output, err
}

//
// CreateDNSRecord creates a new DNS record using the MYRA API
//
func (api *API) CreateDNSRecord(record *DNSRecord, domainName string) (*DNSRecord, error) {
	if _, ok := methods["createDNSRecord"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createDNSRecord")
	}

	definition := methods["createDNSRecord"]
	definition.Action = fmt.Sprintf(definition.Action, domainName)

	result, err := api.call(definition, record)
	if err != nil {
		return nil, err
	}
	return result.(*DNSRecord), nil
}

//
// UpdateDNSRecord updates the passed DNS record using the MYRA API
//
func (api *API) UpdateDNSRecord(record *DNSRecord, domainName string) (*DNSRecord, error) {
	if _, ok := methods["updateDNSRecord"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateDNSRecord")
	}

	definition := methods["updateDNSRecord"]
	definition.Action = fmt.Sprintf(definition.Action, domainName)

	result, err := api.call(definition, record)
	if err != nil {
		return nil, err
	}
	return result.(*DNSRecord), nil
}

//
// DeleteDNSRecord deletes the passed DNS record using the MYRA API
//
func (api *API) DeleteDNSRecord(record *DNSRecord, domainName string) (*DNSRecord, error) {
	if _, ok := methods["deleteDNSRecord"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "deleteDNSRecord")
	}

	definition := methods["deleteDNSRecord"]
	definition.Action = fmt.Sprintf(definition.Action, domainName)

	result, err := api.call(definition, record)
	if err != nil {
		return nil, err
	}
	return result.(*DNSRecord), nil
}
