package documents

// CreateGroupRequest represents the JSON tranport data structure
// for a request to create a Group.
type CreateGroupRequest struct {
	// Schemas is the list of schemas for this API request.
	Schemas []string `json:"schemas"`

	// DisplayName is the human-friendly name given to a group
	// resource.
	DisplayName string `json:"displayName"`
}

// GroupResponse represents the JSON transport data structure
// for a response containing a group resource.
type GroupResponse struct {
	// ID is the unique identifier for a group resource.
	ID string `json:"id"`

	// Schemas is the list of schemas for this API request.
	Schemas []string `json:"schemas"`

	// DisplayName is the human-friendly name given to a group
	// resource.
	DisplayName string `json:"displayName"`

	// Meta is the collection of metadata describing the group
	// resource.
	Meta Meta `json:"meta"`
}

// GroupListResponse represents the JSON transport data structure
// for a response containing a list of group resources.
type GroupListResponse struct {
	// Schemas is the list of schemas for this API request.
	Schemas []string `json:"schemas"`

	// Resources is a list of group resources.
	Resources []GroupResponse `json:"resources"`

	// StartIndex is the index number to start at when returning
	// the list of resources.
	StartIndex int `json:"startIndex"`

	// ItemsPerPage is the number of items to return in the
	// list of resources.
	ItemsPerPage int `json:"itemsPerPage"`

	// TotalResults is the total number of resources that match
	// the list query.
	TotalResults int `json:"totalResults"`
}

// GroupAssociation represents the JSON transport data structure
// for a response contains references to associated groups.
type GroupAssociation struct{}
