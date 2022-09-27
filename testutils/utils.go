package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nsf/jsondiff"
)

// FileToStruct unmarshals a json file into s truct
func FileToStruct(filepath string, s interface{}) io.Reader {
	bb, _ := ioutil.ReadFile(filepath)
	json.Unmarshal(bb, s)
	return bytes.NewReader(bb)
}

// NewTestServer creates a test server for testing
func NewTestServer(h http.HandlerFunc) (string, func()) {
	ts := httptest.NewServer(h)
	return ts.URL, func() { ts.Close() }
}

// GetResponseBody returns the unmarshalled response body
func GetResponseBody(res *http.Response, responseBody interface{}) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(body, responseBody)
}

// AssertResponseBodyEqual asserts that the response body matches the results specified in a file path
func AssertResponseBodyEqual(t *testing.T, resFilePath string, r *http.Response) {
	body, _ := ioutil.ReadAll(r.Body)
	expected, _ := ioutil.ReadFile(resFilePath)
	opts := jsondiff.DefaultConsoleOptions()
	diff, diffStr := jsondiff.Compare(body, expected, &opts)
	if diff != jsondiff.FullMatch {
		t.Errorf(fmt.Sprintf("Diff=%s", diffStr))
	}
}

// SetTestStandardHeaders sets standard headers for testing
func SetTestStandardHeaders(r *http.Request, basicToken string) {
	h := r.Header
	h["Authorization"] = []string{fmt.Sprintf("Bearer %s", basicToken)}
}
