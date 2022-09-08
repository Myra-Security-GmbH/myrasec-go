# File

Using the File API allows basic actions of files within the CDN. These actions are uploading, listing and removing files.

```go
type FileQuery struct {
    Bucket string `json:"bucket,omitempty"`
    Path   string `json:"path,omitempty"`
    Limit  int    `json:"limit,omitempty"`
    Type   int    `json:"type,omitempty"`
    Cursor string `json:"cursor,omitempty"`
}
```

| Field | Type | Description|
|---|---|---|
| `Bucket` | string | The bucket you want to access/query |
| `Path` | string | The path in the bucket you want to access/query |
| `Limit` | int | The maximum amount of elements (files or directories) that will be returned |
| `Type` | int | The type you want to query. `0` for files, `1` for directories |
| `Cursor` | string | The cursor you use - this is required for pagination |


```go
type File struct {
    Type        int             `json:"type"`
    Path        string          `json:"path"`
    Basename    string          `json:"basename"`
    Size        int             `json:"size"`
    Hash        string          `json:"hash"`
    Modified    *types.DateTime `json:"modified"`
    ContentType string          `json:"contentType"`
}
```

| Field | Type | Description|
|---|---|---|
| `Type` | int | The type of the "file". `0` for file, `1` for directory |
| `Path` | string | The path of the file/directory |
| `Basename` | string | The name of the file/directory |
| `Size` | string | The size of the file |
| `Hash` | string | The hash od the file |
| `Modified` | string | Mofification timestamp |
| `ContentType` | string | the content-type of the file |

## Upload files
The path within this request is used to save the file. Sub folders are automatically created, if they are not existing. The content inside the CDN is compressed if it fits some expectations on content and size. For client that are not able to use a compressed format like gzip our nodes will decompress the data for the client before sending.

### Example
```go
file, err := os.Open("/file/to/upload")
if err != nil {
    return err
}
defer file.Close()

err = api.UploadFile(file, "example.com", "b1", "/uploaded/file")
if err != nil {
    return err
}
```

## Upload archive
It is possible to upload a zip archive instead of single files. The API will extract the archive to the given filepath.
The file will be unpacked relative to the given folder which is given by the parameter "path".

### Example
```go
archive, err := os.Open("/file/to/upload.zip")
if err != nil {
    return err
}
defer archive.Close()

err = api.UploadArchive(archive, "example.com", "b1", "/uploaded/archive/")
if err != nil {
    return err
}
```

## List files
You can only list directories (`Type = 1`) or files (`Type = 0`) not together. There is also a limit for the amount of files inside the list at once.
To flip between pages you have to append "cursor", which you will receive as attribute "cursorNext" after the first call.

### Example
```go
query := &myrasec.FileQuery{
    Bucket: "b1",
    Path:   "/",
    Limit:  100,
    Type:   0,
}

cursorNext, files, err := api.ListFiles(query, "example.com")
if err != nil {
    panic(err)
}

log.Println(cursorNext)

for _, f := range files {
    log.Println(f.Basename)
}

nextQuery := &myrasec.FileQuery{
    Bucket: "b1",
    Path:   "/",
    Limit:  100,
    Type:   0,
    Cursor: cursorNext
}
cursorNext, files, err = api.ListFiles(nextQuery, "example.com")
if err != nil {
    panic(err)
}
// ...

```

## Remove files
For deleting a specific file it is necessary to set an absolute path to the file beginning with "/". It is also possible to remove a folder by just pointing the path to a foldername.

### Example
```go
query := &myrasec.FileQuery{
    Bucket: "b1",
    Path:   "/example.jpg",
}

err := api.RemoveFiles(query, "example.com")
if err != nil {
    panic(err)
}
```