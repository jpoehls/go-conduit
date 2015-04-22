package conduit

import (
	_ "github.com/codegangsta/envy/autoload"
	"os"
	"testing"
)

var testHost string
var testUser string
var testCert string

func init() {
	testHost = os.Getenv("GOCONDUIT_HOST")
	testUser = os.Getenv("GOCONDUIT_USER")
	testCert = os.Getenv("GOCONDUIT_CERT")
}

type conduitCapabilitiesResponse struct {
	Authentication []string `json:"authentication"`
	Signatures     []string `json:"signatures"`
	Input          []string `json:"input"`
	Output         []string `json:"output"`
}

func TestCall(t *testing.T) {
	conn, err := Dial(testHost)
	if err != nil {
		t.Fatal(err)
	}

	err = conn.Connect(testUser, testCert)
	if err != nil {
		t.Fatal(err)
	}

	type params struct {
		Names   []string `json:"names"`
		Session *Session `json:"__conduit__"`
	}

	type result map[string]*struct {
		URI      string `json:"uri"`
		FullName string `json:"fullName"`
		Status   string `json:"status"`
	}

	p := &params{
		Names:   []string{"T1"},
		Session: conn.Session,
	}
	var r result

	err = conn.Call("phid.lookup", p, &r)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", r)
}

func TestCapabilities(t *testing.T) {
	var resp conduitCapabilitiesResponse
	err := call(testHost+"/api/conduit.getcapabilities", nil, &resp)

	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v\n", resp)
}

func TestConnect(t *testing.T) {
	conn, err := Dial(testHost)

	if err != nil {
		t.Fatal(err)
	}

	if conn == nil {
		t.Fatal("nil connection")
	}

	err = conn.Connect(testUser, testCert)
	if err != nil {
		t.Fatal(err)
	}

	if conn.host != "https://code.interworks.com" {
		t.Error("bad host")
	}

	if conn.Session == nil {
		t.Error("missing conduit session")
	}

	if conn.Session.SessionKey == "" {
		t.Error("missing session key")
	}

	if conn.Session.ConnectionID == 0 {
		t.Error("missing connectionID")
	}

	r, err := conn.PHIDLookup([]string{"T1", "D1"})
	if err != nil {
		t.Fatal(err)
	}

	if r == nil {
		t.Fatal("nil response")
	}

	t.Logf("%+v\n", r)

	r2, err := conn.PHIDLookupSingle("D1")
	if err != nil {
		t.Fatal(err)
	}
	if r2 == nil {
		t.Fatal("nil response")
	}

	t.Logf("%+v\n", r2)
}
