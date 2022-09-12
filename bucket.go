package myrasec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

// getBucketMethods returns CDN Bucket related API calls
func getBucketMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"listBuckets": {
			BaseURL: "https://upload.myracloud.com/%s",
			Name:    "listBuckets",
			Action:  "v2/bucket/list/%s",
			Method:  http.MethodGet,
			Result:  []Bucket{},
		},
		"getBucketStatus": {
			BaseURL:            "https://upload.myracloud.com/%s",
			Name:               "getBucketStatus",
			Action:             "v2/bucket/status/%s/%s",
			Method:             http.MethodGet,
			Result:             BucketStatus{},
			ResponseDecodeFunc: decodeBucketStatusResponse,
		},
		"getBucketStatistics": {
			BaseURL:            "https://upload.myracloud.com/%s",
			Name:               "getBucketStatistics",
			Action:             "v2/bucket/statistics/%s/%s",
			Method:             http.MethodGet,
			Result:             BucketStatistics{},
			ResponseDecodeFunc: decodeSingleBucketResponse,
		},
		"createBucket": {
			BaseURL:            "https://upload.myracloud.com/%s",
			Name:               "createBucket",
			Action:             "v2/bucket/create/%s",
			Method:             http.MethodPut,
			Result:             Bucket{},
			ResponseDecodeFunc: decodeSingleBucketResponse,
		},
		"linkBucket": {
			BaseURL: "https://upload.myracloud.com/%s",
			Name:    "linkBucket",
			Action:  "v2/bucket/link/%s",
			Method:  http.MethodPut,
			Result:  BucketLink{},
		},
		"unlinkBucket": {
			BaseURL: "https://upload.myracloud.com/%s",
			Name:    "unlinkBucket",
			Action:  "v2/bucket/link/%s",
			Method:  http.MethodDelete,
			Result:  BucketLink{},
		},
		"deleteBucket": {
			BaseURL: "https://upload.myracloud.com/%s",
			Name:    "deleteBucket",
			Action:  "v2/bucket/delete/%s",
			Method:  http.MethodDelete,
			Result:  Bucket{},
		},
	}
}

// bucketResponse ...
type bucketResponse struct {
	Error  bool        `json:"error"`
	Result interface{} `json:"result"`
}

// Bucket ...
type Bucket struct {
	Name          string   `json:"bucket"`
	LinkedDomains []string `json:"linkedDomains"`
}

// BucketStatus ...
type BucketStatus struct {
	Status     string `json:"status"`
	StatusCode int    `json:"statusCode"`
}

// BucketStatistics ...
type BucketStatistics struct {
	Files       int   `json:"files"`
	Folders     int   `json:"folders"`
	StorageSize int64 `json:"storageSize"`
	ContentSize int64 `json:"contentSize"`
}

// BucketLink ...
type BucketLink struct {
	Bucket        string `json:"bucket"`
	SubDomainName string `json:"subDomainName"`
}

// ListBuckets returns a list of all created buckets for the given domain.
func (api *API) ListBuckets(domainName string) ([]Bucket, error) {
	if _, ok := methods["listBuckets"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "listBuckets")
	}

	definition := methods["listBuckets"]
	definition.Action = fmt.Sprintf(definition.Action, domainName)

	result, err := api.call(definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	var records []Bucket
	records = append(records, *result.(*[]Bucket)...)

	return records, nil
}

// GetBucketStatus allows you to get a status for your newly created bucket.
func (api *API) GetBucketStatus(domainName string, bucketName string) (*BucketStatus, error) {
	if _, ok := methods["getBucketStatus"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getBucketStatus")
	}

	definition := methods["getBucketStatus"]
	definition.Action = fmt.Sprintf(definition.Action, domainName, bucketName)

	result, err := api.call(definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	return result.(*BucketStatus), nil
}

// GetBucketStatistics returns statistics for a specific bucket.
func (api *API) GetBucketStatistics(domainName string, bucketName string) (*BucketStatistics, error) {
	if _, ok := methods["getBucketStatistics"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "getBucketStatistics")
	}

	definition := methods["getBucketStatistics"]
	definition.Action = fmt.Sprintf(definition.Action, domainName, bucketName)

	result, err := api.call(definition, map[string]string{})
	if err != nil {
		return nil, err
	}

	return result.(*BucketStatistics), nil
}

// CreateBucket creates a new Bucket for the given domain
func (api *API) CreateBucket(domainName string) (*Bucket, error) {
	if _, ok := methods["createBucket"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "createBucket")
	}

	definition := methods["createBucket"]
	definition.Action = fmt.Sprintf(definition.Action, domainName)

	result, err := api.call(definition, nil)
	if err != nil {
		return nil, err
	}
	return result.(*Bucket), nil
}

// LinkBucket links a sub domain to a bucket
func (api *API) LinkBucket(link *BucketLink, domainName string) (*BucketLink, error) {
	if _, ok := methods["linkBucket"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "linkBucket")
	}

	definition := methods["linkBucket"]
	definition.Action = fmt.Sprintf(definition.Action, domainName)

	_, err := api.call(definition, link)
	if err != nil {
		return nil, err
	}
	return link, nil
}

// UnlinkBucket unlinks a sub domain from a bucket
func (api *API) UnlinkBucket(link *BucketLink, domainName string) (*BucketLink, error) {
	return nil, fmt.Errorf("this action is currently not supported")
}

// DeleteBucket removes a bucket
func (api *API) DeleteBucket(bucket *Bucket, domainName string) (*Bucket, error) {
	if _, ok := methods["deleteBucket"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "deleteBucket")
	}

	definition := methods["deleteBucket"]
	definition.Action = fmt.Sprintf(definition.Action, domainName)

	_, err := api.call(definition, bucket)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

// decodeTagSettingsResponse - custom decode function for bucket status response
func decodeBucketStatusResponse(resp *http.Response, definition APIMethod) (interface{}, error) {
	var res BucketStatus
	err := json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// decodeSingleBucketResponse - custom decode function for bucket responses as they sometimes contain an object instead of an array (result)
func decodeSingleBucketResponse(resp *http.Response, definition APIMethod) (interface{}, error) {
	var res bucketResponse
	err := json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	if res.Error {
		return nil, fmt.Errorf("error in response")
	}

	result := res.Result

	if result == nil {
		return nil, fmt.Errorf("unable to detect element in API response")
	}

	tmp, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	if definition.Result == nil {
		return tmp, nil
	}

	decoder := json.NewDecoder(bytes.NewReader(tmp))
	retValue := reflect.New(reflect.TypeOf(definition.Result))
	retVal := retValue.Interface()
	decoder.Decode(retVal)

	return retVal, err
}
