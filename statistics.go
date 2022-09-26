package myrasec

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// getStatisticMethods returns Statistic related API calls
func getStatisticsMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"queryStatistics": {
			Name:               "queryStatistics",
			Action:             "statistic/query",
			Method:             http.MethodPost,
			Result:             Statistics{},
			ResponseDecodeFunc: decodeStatisticsResponse,
		},
	}
}

// Statistics struct contains the statistical data (Result)
type Statistics struct {
	Query  *StatisticQuery        `json:"query"`
	Result map[string]interface{} `json:"result,omitempty"`
}

// StatisticQuery struct is used to specify the query for the statistical data
type StatisticQuery struct {
	AggregationInterval string                       `json:"aggregationInterval"`
	DataSources         map[string]map[string]string `json:"dataSources"`
	StartDate           *types.DateTime              `json:"startDate"`
	EndDate             *types.DateTime              `json:"endDate"`
	FQDN                []string                     `json:"fqdn"`
	Type                string                       `json:"type"`
}

// QueryStatistics function is used to fetch statistical data
func (api *API) QueryStatistics(query *StatisticQuery) (*Statistics, error) {
	if _, ok := methods["queryStatistics"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "queryStatistics")
	}

	definition := methods["queryStatistics"]

	result, err := api.call(definition, &Statistics{Query: query})
	if err != nil {
		return nil, err
	}
	return result.(*Statistics), nil
}

// decodeStatisticsResponse - custom decode function for statistics response. Used in the QueryStatistics action.
func decodeStatisticsResponse(resp *http.Response, definition APIMethod) (interface{}, error) {
	var res Statistics
	err := json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
