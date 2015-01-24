package documents

type CreateSpaceRequest struct {
	Name             string `json:"name"`
	OrganizationGUID string `json:"organization_guid"`
}
