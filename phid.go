package conduit

type pPHIDLookup struct {
	Names   []string `json:"names"`
	Session *Session `json:"__conduit__"`
}

// PHIDLookupResponse is the result of phid.lookup operations.
type PHIDLookupResponse map[string]*PHIDResult

// PHIDResult is a result item of phid operations.
type PHIDResult struct {
	PHID     string `json:"phid"`
	URI      string `json:"uri"`
	TypeName string `json:"typeName"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	FullName string `json:"fullName"`
	Status   string `json:"status"`
}

// PHIDLookup calls the phid.lookup endpoint.
func (c *Conn) PHIDLookup(names []string) (PHIDLookupResponse, error) {
	p := &pPHIDLookup{
		Names:   names,
		Session: c.Session,
	}

	var r PHIDLookupResponse
	if err := c.Call("phid.lookup", p, &r); err != nil {
		return nil, err
	}

	return r, nil
}

// PHIDLookupSingle calls the phid.lookup endpoint with a single name.
func (c *Conn) PHIDLookupSingle(name string) (*PHIDResult, error) {
	resp, err := c.PHIDLookup([]string{name})
	if err != nil {
		return nil, err
	}

	r := resp[name]
	return r, nil
}

type pPHIDQuery struct {
	PHIDs   []string `json:"phids"`
	Session *Session `json:"__conduit__"`
}

// PHIDQueryResponse is the result of phid.query operations.
type PHIDQueryResponse map[string]*PHIDResult

// PHIDQuery calls the phid.query endpoint.
func (c *Conn) PHIDQuery(phids []string) (PHIDQueryResponse, error) {
	p := &pPHIDQuery{
		PHIDs:   phids,
		Session: c.Session,
	}

	var r PHIDQueryResponse
	if err := c.Call("phid.query", p, &r); err != nil {
		return nil, err
	}

	return r, nil
}

// PHIDQuerySingle calls the phid.query endpoint with a single phid.
func (c *Conn) PHIDQuerySingle(phid string) (*PHIDResult, error) {
	resp, err := c.PHIDQuery([]string{phid})
	if err != nil {
		return nil, err
	}

	r := resp[phid]
	return r, nil
}
