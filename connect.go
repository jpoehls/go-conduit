package conduit

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/karlseguin/typed"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// ConduitError is returned when conduit
// requests return an error response.
type ConduitError struct {
	code string
	info string
}

func (e *ConduitError) Code() string {
	return e.code
}

func (e *ConduitError) Info() string {
	return e.info
}

func (e *ConduitError) Error() string {
	return fmt.Sprintf("%s: %s", e.code, e.info)
}

// IsConduitError checks whether or not e is a ConduitError.
func IsConduitError(e error) bool {
	_, ok := e.(*ConduitError)
	return ok
}

// A Dialer contains options for connecting to an address.
type Dialer struct {
	ClientName        string
	ClientVersion     string
	ClientDescription string
}

// Dial opens a connection to conduit.
func Dial(host, user, cert string) (*Conn, error) {
	var d Dialer
	d.ClientName = "go-conduit"
	d.ClientVersion = "1"
	return d.Dial(host, user, cert)
}

type Conn struct {
	host         string
	sessionKey   string
	connectionID int64
}

func getAuthToken() string {
	return strconv.FormatInt(time.Now().UTC().Unix(), 10)
}

func getAuthSignature(authToken, cert string) string {
	h := sha1.New()
	io.WriteString(h, authToken)
	io.WriteString(h, cert)

	return fmt.Sprintf("%x", h.Sum(nil))
}

type pConduitConnect struct {
	Client            string `json:"client"`
	ClientVersion     string `json:"clientVersion"`
	ClientDescription string `json:"clientDescription"`
	Host              string `json:"host"`
	User              string `json:"user"`
	AuthToken         string `json:"authToken"`
	AuthSignature     string `json:"authSignature"`
}

type conduitConnectResponse struct {
	SessionKey   string `json:"sessionKey"`
	ConnectionID int64  `json:"connectionID"`
}

type conduitAuth struct {
	SessionKey   string `json:"sessionKey"`
	ConnectionID int64  `json:"connectionID"`
	UserPHID     string `json:"userPHID"`
}

func (d *Dialer) Dial(host, user, cert string) (*Conn, error) {
	host = strings.TrimSuffix(host, "/")
	authToken := getAuthToken()
	authSig := getAuthSignature(authToken, cert)

	var resp conduitConnectResponse

	err := call(host+"/api/conduit.connect", &pConduitConnect{
		Client:            d.ClientName,
		ClientVersion:     d.ClientVersion,
		ClientDescription: d.ClientDescription,
		Host:              host,
		User:              user,
		AuthToken:         authToken,
		AuthSignature:     authSig,
	}, &resp)

	if err != nil {
		return nil, err
	}

	conn := Conn{
		host:         host,
		sessionKey:   resp.SessionKey,
		connectionID: resp.ConnectionID,
	}

	return &conn, nil
}

func call(endpointUrl string, params interface{}, resp interface{}) error {
	b, err := json.Marshal(params)
	if err != nil {
		return err
	}

	form := url.Values{}
	form.Add("params", string(b))
	form.Add("output", "json")

	_, isConduitConnect := params.(*pConduitConnect)
	if isConduitConnect {
		form.Add("__conduit__", "true")
	}

	req, err := http.NewRequest("POST", endpointUrl, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	hresp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer hresp.Body.Close()

	body, err := ioutil.ReadAll(hresp.Body)
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

	err = json.Unmarshal(resultBytes, &resp)
	if err != nil {
		return err
	}

	return nil
}
