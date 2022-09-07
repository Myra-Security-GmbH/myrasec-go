# Bucket

```go
type Bucket struct {
    Name          string   `json:"bucket"`
	LinkedDomains []string `json:"linkedDomains"`
}
```

| Field | Type | Description|
|---|---|---|
| `Name` | string | Generated name of a bucket. |
| `LinkedDomains` | []string | List of subdomains, linked to the bucket |

```go
type BucketStatus struct {
	Status     string `json:"status"`
	StatusCode int    `json:"statusCode"`
}
```

| Field | Type | Description|
|---|---|---|
| `Status` | string | Status message |
| `StatusCode` | int | 0 = Bucket available, 1 = Bucket is deleted, 2 = Bucket creating still in progress, 4 = Unknown bucket status |


```go
type BucketStatistics struct {
	Files       int   `json:"files"`
	Folders     int   `json:"folders"`
	StorageSize int64 `json:"storageSize"`
	ContentSize int64 `json:"contentSize"`
}
```

| Field | Type | Description|
|---|---|---|
| `Files` | int | Amount of files |
| `Folders` | int | Amount of folders |
| `StorageSize` | int64 | Size of the bucket |
| `ContentSize` | int64 | Size of the content |


```go
type BucketLink struct {
	Bucket        string `json:"bucket"`
	SubDomainName string `json:"subDomainName"`
}
```

| Field | Type | Description|
|---|---|---|
| `Bucket` | string | The name of the bucket. |
| `SubDomainName` | string | The subdomain name to link to the bucket |


## Create
Create a new bucket for the given domain. The name for a new bucket is generated and will be returned on a successful call.

### Example
```go
b, err := api.CreateBucket("example.com")
if err != nil {
    panic(err)
}

log.Println(b.Name)
```

## Status
Allows you to get a status for your newly created bucket. This API call does not give any information about when your bucket will be created. It only allows you to see if the bucket is ready or not.

### Example
```go
status, err := api.GetBucketStatus("example.com", "b1")
if err != nil {
    panic(err)
}

log.Println(status.Status)
log.Println(status.StatusCode)
```

## Statistics
Returns statistics for a specific bucket. The given sizes are estimations. The sizes can also differ from real content size’s due to the fact that some content might be compressed inside the CDN.

### Example
```go
statistics, err := api.GetBucketStatistics("example.com", "b1")
if err != nil {
    panic(err)
}

if statistics != nil {
    log.Println(statistics.ContentSize)
    log.Println(statistics.Files)
    log.Println(statistics.Folders)
    log.Println(statistics.StorageSize)
}
```

## List
Returns a list of all created buckets for the given domain.

### Example
```go
buckets, err := api.ListBuckets("example.com")
if err != nil {
    panic(err)
}

for _, b := range buckets {
    log.Println(b.Name, b.LinkedDomains)
}
```

## Link
Linking a sub domain to a bucket allows you to publish the content inside the bucket.  
> **You can only link sub domains that are members of the domain given in the URL. The sub domain also needs to be a CDN enabled sub domain in MYRACLOUD to gain access to the published content.**

### Example
```go
link := &myrasec.BucketLink{
    Bucket:        "b1",
    SubDomainName: "www.example.com",
}
b, err := api.LinkBucket(link, "example.com")
if err != nil {
    panic(err)
}
```

## Unlink
Unlinking a sub domain from a bucket. After removing a link you are no longer able to access the bucket’s content for that sub domain.

> **This action is currently not supported. Support for this action will be added in the future.**

## Delete
Removing a bucket will also delete all of the content within the bucket. It is not possible to delete a bucket that have links from sub domains.
> **Removing a bucket deletes also all of it’s content. The content cannot be restored after deletion!**

### Example
```go
bucket := &myrasec.Bucket{
    Name: "b1",
}
b, err := api.DeleteBucket(bucket, "example.com")
if err != nil {
    panic(err)
}
```