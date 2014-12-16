package fakes

type TemplateAssigner struct {
	AssignToClientArguments       []string
	AssignToNotificationArguments []string
	AssignToClientError           error
	AssignToNotificationError     error
}

func NewTemplateAssigner() *TemplateAssigner {
	return &TemplateAssigner{}
}

func (assigner *TemplateAssigner) AssignToClient(clientID, templateID string) error {
	assigner.AssignToClientArguments = []string{clientID, templateID}
	return assigner.AssignToClientError
}

func (assigner *TemplateAssigner) AssignToNotification(clientID, notificationID, templateID string) error {
	assigner.AssignToNotificationArguments = []string{clientID, notificationID, templateID}
	return assigner.AssignToNotificationError
}
