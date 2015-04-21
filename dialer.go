package conduit

import (
	"encoding/json"
	"errors"
	"github.com/karlseguin/typed"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	// ErrJSONOutputUnsupported is returned when conduit doesn't support JSON output.
	ErrJSONOutputUnsupported = errors.New("json output not supported")

	// ErrURLEncodedInputUnsupported is returned when conduit doesn't support URL encoded input.
	ErrURLEncodedInputUnsupported = errors.New("urlencoded input not supported")

	// ErrSessionAuthUnsupported is returned when conduit doesn't support session authentication.
	ErrSessionAuthUnsupported = errors.New("session authentication not supported")
)

// ConduitError is returned when conduit
// requests return an error response.
type ConduitError struct {
	code string
	info string
}

// Code returns the error_code returned in a conduit response.
func (err *ConduitError) Code() string {
	return err.code
}

// Info returns the error_info returned in a conduit response.
func (err *ConduitError) Info() string {
	return err.info
}

func (err *ConduitError) Error() string {
	return err.code + ": " + err.info
}

// IsConduitError checks whether or not err is a ConduitError.
func IsConduitError(err error) bool {
	_, ok := err.(*ConduitError)
	return ok
}

// A Dialer contains options for connecting to an address.
type Dialer struct {
	ClientName        string
	ClientVersion     string
	ClientDescription string
}

// Dial connects to conduit and confirms the API capabilities
// for future calls.
func Dial(host string) (*Conn, error) {
	var d Dialer
	d.ClientName = "go-conduit"
	d.ClientVersion = "1"
	return d.Dial(host)
}

// Dial connects to conduit and confirms the API capabilities
// for future calls.
func (d *Dialer) Dial(host string) (*Conn, error) {
	host = strings.TrimSuffix(host, "/")

	var resp conduitCapabilitiesResponse
	err := call(host+"/api/conduit.getcapabilities", nil, &resp)
	if err != nil {
		return nil, err
	}

	// We use conduit.connect for authentication
	// and it establishes a session.
	if !containsString(resp.Authentication, "session") {
		return nil, ErrSessionAuthUnsupported
	}

	if !containsString(resp.Input, "urlencoded") {
		return nil, ErrURLEncodedInputUnsupported
	}

	if !containsString(resp.Output, "json") {
		return nil, ErrJSONOutputUnsupported
	}

	conn := Conn{
		host:         host,
		capabilities: &resp,
		dialer:       d,
	}

	return &conn, nil
}

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

	if result != nil {
		err = json.Unmarshal(resultBytes, &result)
		if err != nil {
			return err
		}
	}

	return nil
}
