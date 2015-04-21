# Conduit

A Go package for connecting to [Phabricator](http://phabricator.org) via the [Conduit](https://secure.phabricator.com/book/phabdev/article/conduit/) API.

Originally created to help develop [Merlin](https://github.com/jpoehls/merlin), an [Arcanist](http://www.phabricator.com/docs/arcanist/) like tool written in Go.

[Documentation](http://godoc.org/github.com/jpoehls/go-conduit)

# Getting a conduit certificate

This library uses `conduit.connect` to establish an authenticated session. You'll need to have a valid username and conduit certificate in order to use this API.

To get your conduit certificate, go to `https://{MY_PHABRICATOR_URL}/settings/panel/conduit` and copy/paste.

# Usage

## Connecting

```
conn, err := conduit.Dial("https://secure.phabricator.com")

err = conn.Connect("USERNAME", "CERTIFICATE")
```

## Errors

Any conduit error response will be returned as a
`conduit.ConduitError` type

```
conn, err := conduit.Dial("https://secure.phabricator.com")
err = conn.Connect("USERNAME", "CERTIFICATE")

ce, ok := err.(*conduit.ConduitError)
if ok {
	println("code: " + ce.Code())
	println("info: " + ce.Info())
}

// Or
if conduit.IsConduitError(err) {
	// do something else
}
```

## phid.lookup

```
result, err := conduit.PHIDLookup([]string{"T1", "D1"})
```

```
result, err := conduit.PHIDLookupSingle("T1")
```

## phid.query

```
result, err := conduit.PHIDQuery([]string{"PHID-DREV-gumr6ra5wm32ez46qo3f", "..."})
```

```
result, err := conduit.PHIDQuerySingle("PHID-DREV-gumr6ra5wm32ez46qo3f")
```

## Arbitrary calls

You can use the `conn.Call()` method to make arbitrary
conduit method calls that aren't specifically supported
by the package.

```
type params struct {
	Names   []string         `json:"names"`
	Session *conduit.Session `json:"__conduit__"`
}

type result map[string]*struct{
	URI      string `json:"uri"`
	FullName string `json:"fullName"`
	Status   string `json:"status"`
}

p := &params {
	Names: []string{"T1"},
	Session: conn.Session,
}
var r result

err := conn.Call("phid.lookup", p, &r)
```