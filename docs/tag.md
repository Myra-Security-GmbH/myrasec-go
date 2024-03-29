# Tag

```go
type Tag struct {
    ID           int              `json:"id,omitempty"`
    Created      *types.DateTime  `json:"created,omitempty"`
    Modified     *types.DateTime  `json:"modified,omitempty"`
    Name         string           `json:"name"`
    Type         string           `json:"type"`
    Assignments  []TagAssignments `json:"assignments"`
}
```
| Field | Type | Description|
|---|---|---|
| `ID` | int | Id is an unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete a Tag it is necessary to add this attribute to your object. |
| `Created` | *types.DateTime | Created will be created by the server after creating a new Tag object. This value is informational so it is not necessary to add this attribute to any API call. |
| `Modified` | *types.DateTime | Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates and deletes. |
| `Name` | string | Identifies the tag by its name. |
| `Type` | string | Defines the type of the tag and must be one of `CONFIG`, `WAF`, `CACHE`, `RATE_LIMIT` |
| `Assignments` | []TagAssignments |

```go
type TagAssignment struct {
    ID            int             `json:"id,omitempty"`
    Created       *types.DateTime `json:"created,omitempty"`
    Modified      *types.DateTime `json:"modified,omitempty"`
    Type          string          `json:"type"`
    Title         string          `json:"title"`
    SubDomainName string          `json:"subDomainName"`
}
```
| Field | Type | Description|
|---|---|---|
| `ID` | int | Id is an unique identifier for an object. This value is always a number type and cannot be set while inserting a new object. To update or delete a TagAssignment it is necessary to add this attribute to your object. |
| `Created` | *types.DateTime | Created will be created by the server after creating a new TagAssignment object. This value is informational so it is not necessary to add this attribute to any API call. |
| `Modified` | *types.DateTime | Identifies the version of the object. To ensure that you are updating the most recent version and not overwriting other changes, you always have to add modified for updates and deletes. |
| `Type` | string | Defines the type of the tag assignment and must be one of `DOMAIN`, `SUBDOMAIN` |
| `Title` | string | Identifies the tag assignment by its domain name. |
| `SubDomainName` | string | Only set on SUBDOMAIN tag assignments |

## Create
To create a new tag it is neccessary to send a Tag object without the attributes "id" and "modified".
Both attributes will be generated by the server and returned after a successful insert is done.

### Example
```go
newTag := &myrasec.Tag{
    Name: "Example Tag",
    Type: "CONFIG",
    Assignments: []TagAssignment{
        {
            Type: "DOMAIN",
            Title: "example.com",
        }
    }
}

t, err := api.CreateTag(newTag)
if err != nil {
    log.Fatal(err)
}
```

## List
The listing operation returns a list of tags. The list contains tags for the organization of the account you are accessing the API with.

### Example
```go
tags, err := api.ListTags(nil)
if err != nil {
    log.Fatal(err)
}
```

It is possible to pass a map of parameters (`map[string]string`) to the `ListTags` function

| name | description | default |
|---|---|---|
| `search` | Filter by the specified search query | null |
| `page` | Specify the page of the result | 1 |
| `pageSize` | Specify the amount of results in the response | 50 |

## Read
The read operation returns a single Tag by it's ID
```go
tag, err := api.GetTag(tagId)
if err != nil {
    log.Fatal(err)
}
```

## Update
Updating a tag is very similar to creating a new one. The main difference is that an update will need the generated "id" and "modified" attributes to identify the object you are trying to update.

### Example
```go
tag := &myrasec.Tag{
    ID: 0000,
    Modified: &types.DateTime{
        Time: modified
    },
    Name: "Example Tag",
    Type: "CONFIG",
    Assignments: []TagAssignment{
        {
            Type: "DOMAIN",
            Title: "example.com",
        }
    }
}

t, err := api.UpdateTag(tag)
if err != nil {
    log.Fatal(err)
}
```

## Delete
For deleting a tag it is only necessary to send the "id" and "modified" attributes as body content.

### Example
```go
tag := &myrasec.Tag{
    ID: 0000,
    Modified: &types.DateTime{
        Time: modified
    }
}

t, err := api.DeleteTag(tag)
if err != nil {
    log.Fatal(err)
}
```
