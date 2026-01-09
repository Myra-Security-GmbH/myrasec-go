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
	Query  *StatisticQuery `json:"query" jsonschema:"Contains the StatisticQuery."`
	Result map[string]any  `json:"result,omitempty" jsonschema:"Returns the detailed custom data for the requested domain."`
}

// StatisticQuery struct is used to specify the query for the statistical data
type StatisticQuery struct {
	AggregationInterval string                       `json:"aggregationInterval" jsonschema:"The interval for aggregating the data points. The statistics can be requested in various aggregation intervals. The requested data will be split into buckets of the given date interval. This applies only to data requested as histogram. The supported intervals are: 5m, hour, day, and week."`
	DataSources         map[string]map[string]string `json:"dataSources" jsonschema:"List of data sources and output type."`
	StartDate           *types.DateTime              `json:"startDate" jsonschema:"Start of the aggregation interval."`
	EndDate             *types.DateTime              `json:"endDate" jsonschema:"End of the aggregation interval."`
	FQDN                []string                     `json:"fqdn,omitempty" jsonschema:"A list of FQDN."`
	Type                string                       `json:"type" jsonschema:"Mode for selecting domains which should be used."`
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
func decodeStatisticsResponse(resp *http.Response, definition APIMethod) (any, error) {
	var res Statistics
	err := json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
