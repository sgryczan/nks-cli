package extensions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	nks "github.com/NetApp/nks-sdk-go/nks"
)

// SDK is a stub for nks.APIClient. We will define extension methods as an extension of this type, since we
// cannot directly define new methods on imported types.
// See: https://golang.org/doc/effective_go.html#embedding
type SDK struct {
	nks.APIClient
}

// NewClient returns a new api client
func NewClient(token, endpoint string) *SDK {
	c := &SDK{nks.APIClient{
		Token:      token,
		Endpoint:   strings.TrimRight(endpoint, "/"),
		HTTPClient: http.DefaultClient,
	}}
	return c
}

// RunRequest is identical to the original unexported function in the nks sdk
func (c SDK) RunRequest(req *nks.APIReq) error {
	// If method is POST and postObjNeedsEncoding, encode data object and set up payload
	if req.Method == "POST" && req.Payload == nil {
		data, err := json.Marshal(req.PostObj)
		if err != nil {
			return err
		}
		req.Payload = bytes.NewBuffer(data)
	}

	// If path is not fully qualified URL, then prepend with endpoint URL
	if req.Path[0:4] != "http" {
		req.Path = c.Endpoint + req.Path
	}

	// Set up new HTTP request
	httpReq, err := http.NewRequest(req.Method, req.Path, req.Payload)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.Token)
	httpReq.Header.Set("User-Agent", nks.ClientUserAgentString)
	httpReq.Header.Set("Content-Type", "application/json")

	// Run HTTP request, catching response
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return err
	}

	// Check Status Code versus what the caller wanted, error if not correct
	if req.WantedStatus != resp.StatusCode {
		body, _ := ioutil.ReadAll(resp.Body)
		err = fmt.Errorf("Incorrect status code returned: %d, Status: %s\n%s", resp.StatusCode, resp.Status, string(body))
		return err
	}

	// If DELETE operation, return
	if req.Method == "DELETE" || req.ResponseObj == nil {
		return nil
	}

	// Store response from remote server, if not a delete operation
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	req.ResponseString = string(body)

	if req.DontUnmarsahal {
		return err
	}

	// Unmarshal response into ResponseObj struct, return ResponseObj and error, if there is one
	return json.Unmarshal(body, req.ResponseObj)
}
