package myrasec

import (
	"fmt"

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
func (api *API) ListDNSRecords(domainId int, params map[string]string) ([]DNSRecord, error) {
	if _, ok := methods["listDNSRecords"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "listDNSRecords")
	}

	definition := methods["listDNSRecords"]
	definition.Action = fmt.Sprintf(definition.Action, domainId)

	result, err := api.call(definition, params)
	if err != nil {
		return nil, err
	}

	var records []DNSRecord
	for _, v := range *result.(*[]DNSRecord) {
		records = append(records, v)
	}

	return records, nil
}

//
// CreateDNSRecord creates a new DNS record using the MYRA API
//
func (api *API) CreateDNSRecord(record *DNSRecord, domainId int) (*DNSRecord, error) {
	if _, ok := methods["createDNSRecord"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createDNSRecord")
	}

	definition := methods["createDNSRecord"]
	definition.Action = fmt.Sprintf(definition.Action, domainId)

	result, err := api.call(definition, record)
	if err != nil {
		return nil, err
	}
	return result.(*DNSRecord), nil
}

//
// UpdateDNSRecord updates the passed DNS record using the MYRA API
//
func (api *API) UpdateDNSRecord(record *DNSRecord, domainId int) (*DNSRecord, error) {
	if _, ok := methods["updateDNSRecord"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "updateDNSRecord")
	}

	definition := methods["updateDNSRecord"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, record.ID)

	result, err := api.call(definition, record)
	if err != nil {
		return nil, err
	}
	return result.(*DNSRecord), nil
}

//
// DeleteDNSRecord deletes the passed DNS record using the MYRA API
//
func (api *API) DeleteDNSRecord(record *DNSRecord, domainId int) (*DNSRecord, error) {
	if _, ok := methods["deleteDNSRecord"]; !ok {
		return nil, fmt.Errorf("Passed action [%s] is not supported", "deleteDNSRecord")
	}

	definition := methods["deleteDNSRecord"]
	definition.Action = fmt.Sprintf(definition.Action, domainId, record.ID)

	_, err := api.call(definition, record)
	if err != nil {
		return nil, err
	}
	return record, nil
}
