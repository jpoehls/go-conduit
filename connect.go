package conduit

import (
	"crypto/sha1"
	"fmt"
	"io"
	"strconv"
	"time"
)

// Conn is a connection to the conduit API.
type Conn struct {
	host         string
	user         string
	capabilities *conduitCapabilitiesResponse
	Session      *Session
	dialer       *Dialer
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

// Session is the conduit session state
// that will be sent in the JSON params as __conduit__.
type Session struct {
	SessionKey   string `json:"sessionKey"`
	ConnectionID int64  `json:"connectionID"`
}

// Connect calls conduit.connect to open an authenticated
// session for future requests.
func (c *Conn) Connect(user, cert string) error {
	authToken := getAuthToken()
	authSig := getAuthSignature(authToken, cert)
	c.user = user

	var resp conduitConnectResponse

	if err := c.Call("conduit.connect", &pConduitConnect{
		Client:            c.dialer.ClientName,
		ClientVersion:     c.dialer.ClientVersion,
		ClientDescription: c.dialer.ClientDescription,
		Host:              c.host,
		User:              c.user,
		AuthToken:         authToken,
		AuthSignature:     authSig,
	}, &resp); err != nil {
		return err
	}

	c.Session = &Session{
		SessionKey:   resp.SessionKey,
		ConnectionID: resp.ConnectionID,
	}

	return nil
}

// Call allows you to make a raw conduit method call.
// Params will be marshalled as JSON and the result JSON
// will be unmarshalled into the result interface{}.
//
// This is primarily useful for calling conduit endpoints that
// aren't specifically supported by other methods in this
// package.
func (c *Conn) Call(method string, params interface{}, result interface{}) error {
	err := call(c.host+"/api/"+method, params, &result)
	return err
}
