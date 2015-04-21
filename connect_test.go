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

	if conn.conduitAuth == nil {
		t.Error("missing conduit auth")
	}

	if conn.conduitAuth.SessionKey == "" {
		t.Error("missing session key")
	}

	if conn.conduitAuth.ConnectionID == 0 {
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

/* Test Helpers */
// func expect(t *testing.T, a interface{}, b interface{}) {
// 	if a != b {
// 		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
// 	}
// }
