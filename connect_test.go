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

func TestConnect(t *testing.T) {
	conn, err := Dial(testHost, testUser, testCert)

	if err != nil {
		t.Fatal(err)
	}

	if conn == nil {
		t.Fatal("nil connection")
	}

	if conn.host != "https://code.interworks.com" {
		t.Error("bad host")
	}

	if conn.sessionKey == "" {
		t.Error("missing sessionKey")
	}

	if conn.connectionID == 0 {
		t.Error("missing connectionID")
	}

	r, err := conn.PhidLookup([]string{"T1", "D1"})
	if err != nil {
		t.Fatal(err)
	}

	if r == nil {
		t.Fatal("nil response")
	}

	t.Logf("%+v\n", r)

	r2, err := conn.PhidLookupSingle("D1")
	if err != nil {
		t.Fatal(err)
	}
	if r2 == nil {
		t.Fatal("nil response")
	}

	t.Logf("%+v\n", r2)
}
