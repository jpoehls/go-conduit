package conduit

import (
	"testing"
)

func TestCreate(t *testing.T) {

	conn, err := Dial(testHost)
	if err != nil {
		t.Fatal(err)
	}

	err = conn.Connect(testUser, testCert)
	if err != nil {
		t.Fatal(err)
	}

	params := &PasteCreateParams{
		Content:  "test content",
		Title:    "test title",
		Language: "perl",
	}

	item, err := conn.PasteCreate(params)
	if err != nil {
		t.Fatal(err)
	}

	if item == nil {
		t.Fatal("returned paste is nil")
	}

	expect(t, item.Content, "test content")
	expect(t, item.Title, "test title")
	expect(t, item.Language, "perl")

	t.Logf("%+v\n", item)
}
