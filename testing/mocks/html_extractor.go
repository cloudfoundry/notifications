package mocks

type HTMLExtractor struct {
	ExtractCall struct {
		Returns struct {
			Error error
		}
	}
}

func NewHTMLExtractor() *HTMLExtractor {
	return &HTMLExtractor{}
}

func (e HTMLExtractor) Extract(html string) (doctype, head, bodyContent, bodyAttributes string, err error) {
	return "", "", "", "", e.ExtractCall.Returns.Error
}
