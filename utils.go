package conduit

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/karlseguin/typed"
)

// containsString checks whether s contains e.
func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// call does the heavy lifting of sending a request to conduit,
// handling error responses by returning *ConduitError,
// and unmarshalling the JSON result into the specified
// result interface{}.
func call(endpointURL string, params interface{}, result interface{}) error {
	form := url.Values{}
	form.Add("output", "json")

	if params != nil {
		b, err := json.Marshal(params)
		if err != nil {
			return err
		}

		form.Add("params", string(b))

		_, isConduitConnect := params.(*pConduitConnect)
		if isConduitConnect {
			form.Add("__conduit__", "true")
		}
	}

	req, err := http.NewRequest("POST", endpointURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	jsonBody, err := typed.Json(body)
	if err != nil {
		return err
	}

	// parse any error conduit returned first
	if jsonBody.String("error_code") != "" {
		return &ConduitError{
			code: jsonBody.String("error_code"),
			info: jsonBody.String("error_info"),
		}
	}

	// if no error, parse the expected result
	resultBytes, err := jsonBody.ToBytes("result")
	if err != nil {
		return err
	}

	if result != nil && resultBytes != nil {
		if err = json.Unmarshal(resultBytes, &result); err != nil {
			return err
		}
	}

	return nil
}

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}
