package warrant

import "time"

// Group is the representation of a group resource within UAA.
type Group struct {
	// ID is the unique identifier for the group resource.
	ID string

	// DisplayName is the human-friendly name given to a group.
	DisplayName string

	// Version is an integer value indicating which revision this resource represents.
	Version int

	// CreatedAt is a timestamp value indicating when the group was created.
	CreatedAt time.Time

	// UpdatedAt is a timestamp value indicating when the group was last modified.
	UpdatedAt time.Time
}
