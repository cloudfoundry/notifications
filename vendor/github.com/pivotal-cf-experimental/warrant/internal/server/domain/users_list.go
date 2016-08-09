package domain

import "github.com/pivotal-cf-experimental/warrant/internal/documents"

type UsersList []User

func (ul UsersList) ToDocument() documents.UserListResponse {
	doc := documents.UserListResponse{
		ItemsPerPage: 100,
		StartIndex:   1,
		TotalResults: len(ul),
		Schemas:      schemas,
		Resources:    []documents.UserResponse{},
	}

	for _, user := range ul {
		doc.Resources = append(doc.Resources, user.ToDocument())
	}

	return doc
}
