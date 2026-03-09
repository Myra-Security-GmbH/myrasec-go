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

// Statistics acts as a container for the requested statistical report.
// It includes both the original query parameters and the resulting data set.
type Statistics struct {
	// Query echoes the parameters used to generate this statistical report.
	Query *StatisticQuery `json:"query" jsonschema:"The query configuration used to generate the results."`

	// Result contains the raw statistical data.
	// The structure of this map depends dynamically on the requested 'dataSources'.
	Result map[string]any `json:"result,omitempty" jsonschema:"The detailed statistical data payload. The keys correspond to the requested data sources (e.g., 'stats', 'waf'). Structure varies based on the query."`
}

// StatisticQuery defines the criteria for fetching statistical data.
// It specifies the time range, target domains, and specific metrics (data sources) to retrieve.
type StatisticQuery struct {
	// AggregationInterval defines the time granularity for histogram data.
	// Only relevant if the data source type is 'histogram'.
	AggregationInterval string `json:"aggregationInterval" jsonschema:"The time bucket size for aggregating data points. Valid values: '5m' (5 minutes), 'hour', 'day', 'week'. Only applies to 'histogram' output types."`

	// DataSources configures which metrics to fetch and how to format them.
	// Map Key: The metric source name (e.g., 'stats', 'waf', 'dns', 'cache').
	// Map Value: Configuration options (e.g., {'type': 'histogram'} or {'type': 'top', 'limit': '10'}).
	DataSources map[string]map[string]string `json:"dataSources" jsonschema:"Specifies the metrics to retrieve. Keys are source names (e.g., 'stats', 'waf'); values are configuration maps (e.g., 'type'='histogram')."`

	// StartDate is the beginning of the reporting window (ISO 8601).
	StartDate *types.DateTime `json:"startDate" jsonschema:"The start timestamp (ISO 8601) for the data aggregation interval."`

	// EndDate is the end of the reporting window (ISO 8601).
	EndDate *types.DateTime `json:"endDate" jsonschema:"The end timestamp (ISO 8601) for the data aggregation interval."`

	// FQDN is a list of specific Fully Qualified Domain Names to include.
	// Only used if Type is set to 'fqdn'.
	FQDN []string `json:"fqdn,omitempty" jsonschema:"List of specific Fully Qualified Domain Names (FQDNs) to filter by. Required if 'type' is set to 'fqdn'."`

	// Type defines the scope of the query.
	// Common values are 'fqdn' (specific domains).
	Type string `json:"type" jsonschema:"The scope selection mode. Valid values: 'fqdn' (use the provided FQDN list)."`
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
	res, ok := result.(*Statistics)
	if !ok {
		return nil, fmt.Errorf("unexpected result type %T", result)
	}
	return res, nil
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
