# Statistics
The statistics API lets you fetch statistical data of your domains like requests, traffic, performance, or
health.

```go
type Statistics struct {
	Query  *StatisticQuery        `json:"query"`
	Result map[string]interface{} `json:"result,omitempty"`
}
```

| Field | Type | Description|
|---|---|---|
| `Query` | *StatisticQuery | Contains the StatisticQuery |
| `Result` | map[string]interface | Returns the detailed custom data for the requested domain |


```go
type StatisticQuery struct {
	AggregationInterval string                       `json:"aggregationInterval"`
	DataSources         map[string]map[string]string `json:"dataSources"`
	StartDate           *types.DateTime              `json:"startDate"`
	EndDate             *types.DateTime              `json:"endDate"`
	FQDN                []string                     `json:"fqdn"`
	Type                string                       `json:"type"`
}
```

| Field | Type | Description|
|---|---|---|
| `AggregationInterval` | string | The interval for aggregating the data points |
| `DataSources` | map[string]map[string]string | List of data sources and output type |
| `StartDate` | *types.DateTime | Start of the aggregation interval |
| `EndDate` | *types.DateTime | End of the aggregation interval |
| `FQDN` | []string | A list of FQDN |
| `Type` | string | Mode for selecting domains which should be used |

### AggregationInterval
The statistics can be requested in various aggregation intervals. The requested data will be split into
buckets of the given date interval. This applies only to data requested as histogram. The supported
intervals are: "5m", "hour", "day", and "week".  

### DataSources
With the statistics API, you can query various information about request types and how they were handled.
```json
{
    "name" : {
        "source" : "bytes_cached",
        "type" : "stats"
    }
}
```
> | Attribute | Type | Description|
> |---|---|---|
> | `name` | string | Arbitrary name of the dataset |
> | `source` | string | Category of the requested data source |
> | `type` | string | Type of the data aggregation |

> ### Name
> The given name is used to name the corresponding result set in the API response. The name may
only contain [a-zA-Z0-9_] characters.
> ### Source
> See the tables [Request data sources](#request-data-sources) and [Traffic data sources](#traffic-data-sources) for a list of possible data source names
to use.
> ### Type
> The statistic data can be requested in two different aggregation types.
>> **stats**: : Data as object containing min/max/avg/sum values.  
>> **histogram**: : Response will consist of multiple objects containing the value for every aggregation bucket.

### StartDate
Included start of the requested period.

### EndDate
Included end of the requested period.

### FQDN
Contains a list of FQDN for which statistics should be generated. Note that you can also use *`ALL:fqdn.de`* as domain name to include data for all subdomains. This value is only used if the `Type` is set to ’fqdn’.

### Type
Mode for selecting domains which should be used.  
**fqdn**: Process the FQDN list given in `FQDN`


### Request data sources
Myra distinguishes incoming requests as SSL and non-SSL depending on the protocol used by the client initiating the request. You can also retrieve information about whether the response was sent from the Myra cache or from origin system.

| Data source | S | N | C | U | Description|
|---|---|---|---|---|---|
| `requests` | X | X | X | X | Total amount of requests handled by Myra |
| `requests_ssl` | X |   | X | X | Amount of requests received via SSL |
| `requests_nonssl` |   | X | X | X | Amount of requests received not via SSL |
| `requests_cached` | X | X | X |   | Total amount of requests delivered from cache |
| `requests_cached_ssl` | X |   | X |   | Amount of requests via SSL delivered from cache |
| `requests_cached_nonssl` |   | X | X |   | Amount of requests via non-SSL delivered from cache |
| `requests_uncached` | X | X |   | X | Total amount of requests passed to the origin |
| `requests_uncached_ssl` | X |   |   | X | Total amount of requests via SSL passed to the origin |
| `requests_uncached_nonssl` |   | X |   | X | Total amount of requests via non-SSL passed to the origin |
| `requests_cache_hits` | X | X | X |   | Ratio of total cached requests to uncached requests in percent |
| `requests_cache_hits_ssl` | X |   | X |   | Ratio of SSL cached requests to SSL uncached requests in percent |
| `requests_cache_hits_nonssl` |   | X | X |   | Ratio of non-SSL cached requests to non-SSL uncached requests in percent |

**S** = SSL, **N** = non-SSL, **C** = Cached, **U** = Uncached

> **Data source**: Name of the data source.  
> **SSL**: Includes amount of requests sent and answered via SSL connection.  
> **non-SSL**: Includes amount of requests sent and answered via non-SSL connection.  
> **Cached**: Includes amount of requests answered from the Myra cache without querying the origin server.  
> **Uncached**: Includes amount of requests answered by Myra after passing the request to the origin server and sending the corresponding response.

### Traffic data sources 
Myra distinguishes transferred data as SSL and non-SSL traffic depending on the protocol used by the client initiating the request. You can also retrieve information about whether the response was sent from the Myra cache or was fetched from the origin system.

| Data source | S | N | C | U | Description|
|---|---|---|---|---|---|
| `bytes` | X | X | X | X | Total amount of outgoing data in bytes |
| `bytes_ssl` | X |   | X | X | Amount of data sent via SSL |
| `bytes_nonssl` |   | X | X | X | Amount of data sent via non-SSL |
| `bytes_cached` | X | X | X |   | Total amount of data delivered from the cache |
| `bytes_cached_ssl` | X |   | X |   | Amount of data via SSL delivered from the cache |
| `bytes_cached_nonssl` |   | X | X |   | Amount of data via non-SSL delivered from the cache |
| `bytes_uncached` | X | X |   | X | Total amount of data passed through from the origin |
| `bytes_uncached_ssl` | X |   |   | X | Amount of data passed via SSL from the origin |
| `bytes_uncached_nonssl` |   | X |   | X | Amount of data passed via non-SSL from the origin |
| `bytes_cache_hits` | X | X | X |   | Ratio of total bytes delivered from the cache in percent |
| `bytes_cache_hits_ssl` | X |   | X |   | Ratio of bytes delivered via SSL from the cache in percent |
| `bytes_cache_hits_nonssl` |   | X | X |   | Ratio of bytes delivered via non-SSL from the cache in percent |

**S** = SSL, **N** = non-SSL, **C** = Cached, **U** = Uncached

> **Data source**: Name of the data source.  
> **SSL**: Includes bytes transferred via SSL connection.  
> **non-SSL**: Includes bytes transferred via non-SSL connection.  
> **Cached**: Includes bytes of all responses transferred from Myra cache without querying the origin server.  
> **Uncached**: Includes bytes of all responses transferred from Myra without caching.

### Other
In addition to traffic and request statistics, Myra allows you to view detailled information about geographic distribution of your visitors, the performance of your origin server and the HTTP response codes of your application.

| Data source | S | N | C | U | Description|
|---|---|---|---|---|---|
| `upstream_performance` | X | X |   | X | Average upstream response time |
| `response_codes` | X | X | X | X | HTTP response codes of total requests |
| `country_codes` | X | X | X | X | Total requests by country |

**S** = SSL, **N** = non-SSL, **C** = Cached, **U** = Uncached

### Result
```json
{
    "bytes_cached_stats": {
        "avg" : 5866621.8906448,
        "max" : 65366760,
        "min" : 0,
        "sum" : 18561991662
    }
    "requests_histogram": {
        "1422399600000": {
            "avg" : 334.70860927152,
            "max" : 1609,
            "min" : 1,
            "sum" : 101082
        }
    }
}
```

| Attribute | Type | Description|
|---|---|---|
| `name` | string | Name of the result set according to chosen name in the request |
| `avg` | Float | The average value of the requested source |
| `max` | Float | The maximal value of the requested source |
| `min` | Float | The minimal value of the requested source |
| `sum` | Float | The sum of the requested source |

**Info**  
The structure of the returned data depends on the requested type. Data sources requested with type "stats" will be returned as an object with min/max/avg/sum keys as seen in "bytes_cached_stats" above. Responses for type "histogram" will contain multiple objects with the key "value", one object for every aggregation bucket. Please note that the `timestamp` used as key is in `milliseconds`.

## Query

### Example
```go
query := &myrasec.StatisticQuery{
    StartDate:           &types.DateTime{Time: time.Now().Add(-time.Hour * 24)},
    EndDate:             &types.DateTime{Time: time.Now()},
    FQDN:                []string{"www.example.com"},
    AggregationInterval: "hour",
    Type:                "fqdn",
    DataSources: map[string]map[string]string{
        "requests_histogram": {
            "source": "requests",
            "type":   "histogram",
        },
        "url_hits_top": {
            "source": "url_hits",
            "type":   "top",
        },
        "cache_misses_top": {
            "source": "cache_misses",
            "type":   "top",
        },
    },
}

statistics, err := api.QueryStatistics(query)
if err != nil {
    log.Fatal(err)
}

for k, v := range statistics.Result {
    log.Println(k, v)
}
```

