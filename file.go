package myrasec

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Myra-Security-GmbH/myrasec-go/v2/pkg/types"
)

// getFileMethods returns CDN File related API calls
func getFileMethods() map[string]APIMethod {
	return map[string]APIMethod{
		"uploadFile": {
			BaseURL: "https://upload.myracloud.com/%s",
			Name:    "uploadFile",
			Action:  "v2/upload/%s/%s/%s",
			Method:  http.MethodPut,
		},
		"listFiles": {
			BaseURL:            "https://upload.myracloud.com/%s",
			Name:               "listFiles",
			Action:             "v2/list/%s",
			Method:             http.MethodPost,
			Result:             listFilesResponse{},
			ResponseDecodeFunc: decodeListFilesResponse,
		},
		"removeFiles": {
			BaseURL: "https://upload.myracloud.com/%s",
			Name:    "removeFiles",
			Action:  "v2/delete/%s",
			Method:  http.MethodDelete,
		},
	}
}

// listFilesResponse ...
type listFilesResponse struct {
	Error      bool   `json:"error"`
	CursorNext string `json:"cursorNext"`
	Files      []File `json:"result"`
	Count      int    `json:"count"`
}

// FileQuery ...
type FileQuery struct {
	Bucket string `json:"bucket,omitempty"`
	Path   string `json:"path,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Type   int    `json:"type,omitempty"`
	Cursor string `json:"cursor,omitempty"`
}

// File ...
type File struct {
	Type        int             `json:"type"`
	Path        string          `json:"path"`
	Basename    string          `json:"basename"`
	Size        int             `json:"size"`
	Hash        string          `json:"hash"`
	Modified    *types.DateTime `json:"modified"`
	ContentType string          `json:"contentType"`
}

// UploadFile ...
func (api *API) UploadFile(domainName string, bucketName string, path string) error {
	return nil
}

// ListFiles returns a list of files or directories as specified in the *FileQuery parameter
func (api *API) ListFiles(fileQuery *FileQuery, domainName string) (string, []File, error) {
	if _, ok := methods["listFiles"]; !ok {
		return "", nil, fmt.Errorf("passed action [%s] is not supported", "listFiles")
	}

	definition := methods["listFiles"]
	definition.Action = fmt.Sprintf(definition.Action, domainName)

	result, err := api.call(definition, fileQuery)
	if err != nil {
		return "", nil, err
	}

	resp := result.(*listFilesResponse)

	return resp.CursorNext, resp.Files, nil
}

// RemoveFiles removes the file or directoy as specified in the *FileQuery param
func (api *API) RemoveFiles(fileQuery *FileQuery, domainName string) error {
	if _, ok := methods["removeFiles"]; !ok {
		return fmt.Errorf("passed action [%s] is not supported", "removeFiles")
	}

	definition := methods["removeFiles"]
	definition.Action = fmt.Sprintf(definition.Action, domainName)

	_, err := api.call(definition, fileQuery)
	if err != nil {
		return err
	}

	return nil
}

// decodeListFilesResponse ...
func decodeListFilesResponse(resp *http.Response, definition APIMethod) (interface{}, error) {
	var res listFilesResponse
	err := json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
