package conduit

type pPhidLookup struct {
	Names   []string     `json:"names"`
	Conduit *conduitAuth `json:"__conduit__"`
}

type PhidLookupResponse map[string]*PhidLookupResult

type PhidLookupResult struct {
	PHID     string `json:"phid"`
	Uri      string `json:"uri"`
	TypeName string `json:"typeName"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	FullName string `json:"fullName"`
	Status   string `json:"status"`
}

func (c *Conn) PhidLookup(names []string) (PhidLookupResponse, error) {
	p := &pPhidLookup{
		Names: names,
		Conduit: &conduitAuth{
			SessionKey:   c.sessionKey,
			ConnectionID: c.connectionID,
		},
	}

	var r PhidLookupResponse
	err := call(c.host+"/api/phid.lookup", p, &r)

	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c *Conn) PhidLookupSingle(name string) (*PhidLookupResult, error) {
	resp, err := c.PhidLookup([]string{name})
	if err != nil {
		return nil, err
	}

	r := resp[name]
	return r, nil
}
