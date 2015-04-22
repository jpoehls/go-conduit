package conduit

import (
	"errors"
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

type conduitCapabilitiesResponse struct {
	Authentication []string `json:"authentication"`
	Signatures     []string `json:"signatures"`
	Input          []string `json:"input"`
	Output         []string `json:"output"`
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
	if err := call(host+"/api/conduit.getcapabilities", nil, &resp); err != nil {
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
