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
	conduitAuth  *conduitAuth
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

type conduitAuth struct {
	SessionKey   string `json:"sessionKey"`
	ConnectionID int64  `json:"connectionID"`
	UserPHID     string `json:"userPHID"`
}

// Connect calls conduit.connect to open an authenticated
// session for future requests.
func (c *Conn) Connect(user, cert string) error {
	authToken := getAuthToken()
	authSig := getAuthSignature(authToken, cert)
	c.user = user

	var resp conduitConnectResponse

	err := call(c.host+"/api/conduit.connect", &pConduitConnect{
		Client:            c.dialer.ClientName,
		ClientVersion:     c.dialer.ClientVersion,
		ClientDescription: c.dialer.ClientDescription,
		Host:              c.host,
		User:              c.user,
		AuthToken:         authToken,
		AuthSignature:     authSig,
	}, &resp)

	if err != nil {
		return err
	}

	c.conduitAuth = &conduitAuth{
		SessionKey:   resp.SessionKey,
		ConnectionID: resp.ConnectionID,
	}

	return nil
}
