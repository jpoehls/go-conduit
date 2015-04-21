# Conduit

A Go package for connecting to Phabricator via the Conduit API.

[Documentation](http://godoc.org/github.com/jpoehls/go-conduit)

## Usage

### Connecting

```
conn, err := conduit.Dial("https://secure.phabricator.com")

err = conn.Connect("USERNAME", "CERTIFICATE")
```

### Errors

Any conduit error response will be returned as a
`conduit.ConduitError` type

```
_, err := conduit.Dial("bad url", "bad user", "bad cert")
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

### phid.lookup

```
result, err := conduit.PHIDLookup([]string{"T1", "D1"})
```

```
result, err := conduit.PHIDLookupSingle("T1")
```

### phid.query

```
result, err := conduit.PHIDQuery([]string{"PHID-DREV-gumr6ra5wm32ez46qo3f", "PHID-DREV-gumr6ra5wm32ez483827"})
```

```
result, err := conduit.PHIDQuerySingle("PHID-DREV-gumr6ra5wm32ez46qo3f")
```