package fakes

type TemplateAssigner struct {
	AssignToClientArguments []string
	AssignToClientError     error
}

func NewTemplateAssigner() *TemplateAssigner {
	return &TemplateAssigner{}
}

func (assigner *TemplateAssigner) AssignToClient(clientID, templateID string) error {
	assigner.AssignToClientArguments = []string{clientID, templateID}
	return assigner.AssignToClientError
}
