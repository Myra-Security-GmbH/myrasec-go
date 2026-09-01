package myrasec

import (
	"bytes"
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
// It specifies the time range, target domains and specific metrics (data sources) to retrieve.
type StatisticQuery struct {
	// AggregationInterval defines the time granularity for histogram data.
	// Only relevant if the data source type is 'histogram'. Omitted = the API default ('day').
	AggregationInterval string `json:"aggregationInterval,omitempty" jsonschema:"The time bucket size for aggregating data points. Valid values: '2m' (2 minutes), 'hour', 'day' (default), 'week', 'month', 'year'. Only applies to 'histogram' output types."`

	// DataSources configures which metrics to fetch and how to format them.
	// Map Key: A result set name of your choice ([a-zA-Z0-9_]); the API returns the data under this name.
	// Map Value: {"source": "<data source>", "type": "<stats|histogram|top>"}, e.g. {"source": "bytes", "type": "stats"}.
	// See docs/statistics.md for the list of data sources.
	DataSources map[string]map[string]string `json:"dataSources" jsonschema:"Result sets to compute, keyed by a name of your choice ([a-zA-Z0-9_]). Each value is a map with 'source' (data source name, e.g. 'requests', 'bytes', 'requests_blocked', 'url_hits') and 'type' ('stats' = min/max/avg/sum, 'histogram' = one value per aggregationInterval bucket, 'top' = ranked URL list)."`

	// StartDate is the beginning of the reporting window (ISO 8601).
	StartDate *types.DateTime `json:"startDate" jsonschema:"The start timestamp (ISO 8601) for the data aggregation interval."`

	// EndDate is the end of the reporting window (ISO 8601). The API accepts at most 7952399 seconds
	// (about 92 days) between StartDate and EndDate.
	EndDate *types.DateTime `json:"endDate" jsonschema:"The end timestamp (ISO 8601) for the data aggregation interval; at most 92 days after startDate."`

	// FQDN is a list of specific Fully Qualified Domain Names to include (at most 150).
	// Only used if Type is set to 'fqdn'. 'ALL:example.com' includes the domain with all its subdomains,
	// but the API resets the list when it expands such an entry: every FQDN listed BEFORE it is discarded.
	// Put an ALL: entry first and use at most one. Results are aggregated over all listed FQDNs.
	FQDN []string `json:"fqdn,omitempty" jsonschema:"List of specific Fully Qualified Domain Names (FQDNs) to filter by, at most 150. 'ALL:example.com' = the domain with all its subdomains; the API discards every FQDN listed before such an entry, so put it first and use at most one. Results are aggregated over all listed FQDNs. Required if 'type' is set to 'fqdn'."`

	// Type defines the scope of the query. Omitted = the API default ('fqdn').
	// 'fqdn' uses the FQDN list; 'all', 'own' and 'foreign' select the domains the user can access,
	// owns, or has been assigned (subdomains included, FQDN ignored). The API rejects these three
	// when they expand to more than 150 domains.
	Type string `json:"type,omitempty" jsonschema:"The scope selection mode. Valid values: 'fqdn' (default, use the provided FQDN list), 'all', 'own' or 'foreign' (every accessible, owned or assigned domain with its subdomains; rejected above 150 domains)."`
}

// QueryStatistics function is used to fetch statistical data
func (api *API) QueryStatistics(query *StatisticQuery) (*Statistics, error) {
	if _, ok := api.methods["queryStatistics"]; !ok {
		return nil, fmt.Errorf("passed action [%s] is not supported", "queryStatistics")
	}

	definition := api.methods["queryStatistics"]

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

// statisticsEnvelope is the wire shape of a statistic/query response. The API
// answers validation failures (unknown or denied FQDN, malformed dates, too
// many FQDNs) with HTTP 200 and the usual error envelope, so the error fields
// must be decoded alongside the payload; the generic Response type cannot be
// used because its "result" field is a list while the statistics result is
// an object. Result stays raw because an empty result set is serialized as
// [] (empty PHP array), which must not fail the decode.
type statisticsEnvelope struct {
	Error         bool            `json:"error"`
	ErrorMessage  string          `json:"errorMessage"`
	ViolationList []*Violation    `json:"violationList"`
	Query         *StatisticQuery `json:"query"`
	Result        json.RawMessage `json:"result"`
}

// decodeStatisticsResponse - custom decode function for statistics response. Used in the QueryStatistics action.
// An error envelope (error: true) is returned as a Go error carrying the violations, never as an empty result.
func decodeStatisticsResponse(resp *http.Response, definition APIMethod) (any, error) {
	var envelope statisticsEnvelope
	err := json.NewDecoder(resp.Body).Decode(&envelope)
	if err != nil {
		return nil, err
	}
	if envelope.Error {
		return nil, formatAPIError(envelope.ErrorMessage, envelope.ViolationList)
	}
	stats := &Statistics{Query: envelope.Query}
	if raw := bytes.TrimSpace(envelope.Result); len(raw) > 0 && raw[0] == '{' {
		if err := json.Unmarshal(raw, &stats.Result); err != nil {
			return nil, err
		}
	}
	return stats, nil
}
