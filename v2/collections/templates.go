package collections

type Template struct {
	ID       string
	Name     string
	HTML     string
	Text     string
	Subject  string
	Metadata string
}

type TemplatesCollection struct {
}

func NewTemplatesCollection() TemplatesCollection {
	return TemplatesCollection{}
}

func (c TemplatesCollection) Get(conn ConnectionInterface, templateID, clientID string) (Template, error) {
	panic("not implemented")
}
