package conduit

// PasteCreateParams are the parameters
// for PasteCreate.
type PasteCreateParams struct {
	Content  string `json:"content"`  // required
	Title    string `json:"title"`    // optional
	Language string `json:"language"` // optional
}

type pPasteCreate struct {
	PasteCreateParams
	Session *Session `json:"__conduit__"`
}

// PasteItem is a result item for paste queries.
type PasteItem struct {
	ID         uint64 `json:"id"`
	ObjectName string `json:"objectName"`
	PHID       string `json:"phid"`
	AuthorPHID string `json:"authorPHID"`
	FilePHID   string `json:"filePHID"`
	Title      string `json:"title"`

	// TODO: figure out how to marshall/unmarshal the unix time
	//       that this comes in
	//DateCreated time.Time `json:"dateCreated"`

	Language   string `json:"language"`
	URI        string `json:"uri"`
	ParentPHID string `json:"parentPHID"`
	Content    string `json:"content"`
}

// PasteCreate calls the paste.create endpoint.
func (c *Conn) PasteCreate(params *PasteCreateParams) (*PasteItem, error) {
	p := &pPasteCreate{
		Session: c.Session,
	}
	p.Content = params.Content
	p.Title = params.Title
	p.Language = params.Language

	var r PasteItem
	if err := c.Call("paste.create", p, &r); err != nil {
		return nil, err
	}

	return &r, nil
}

// PasteQueryParams are the parameters
// for PasteQuery.
type PasteQueryParams struct {
	IDs         []uint64 `json:"ids"`         // optional
	PHIDs       []string `json:"phids"`       // optional
	AuthorPHIDs []string `json:"authorPHIDs"` // optional
	Offset      uint64   `json:"after"`       // optional
	Limit       uint64   `json:"limit"`       // optional
}

type pPasteQuery struct {
	PasteQueryParams
	Session *Session `json:"__conduit__"`
}

// PasteQueryResponse is the result of PasteQuery.
type PasteQueryResponse []*PasteItem

// PasteQuery calls the paste.query endpoint.
func (c *Conn) PasteQuery(params *PasteQueryParams) (PasteQueryResponse, error) {
	p := &pPasteQuery{
		Session: c.Session,
	}
	p.IDs = params.IDs
	p.PHIDs = params.PHIDs
	p.AuthorPHIDs = params.AuthorPHIDs
	p.Offset = params.Offset
	p.Limit = params.Limit

	var items PasteQueryResponse

	var r map[string]*PasteItem
	if err := c.Call("paste.query", p, &r); err != nil {
		return nil, err
	}

	for _, v := range r {
		items = append(items, v)
	}

	return items, nil
}
