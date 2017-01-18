package domain

import "github.com/pivotal-cf-experimental/warrant/internal/documents"

type ClientsList []Client

func (cl ClientsList) ToDocument() documents.ClientListResponse {
	doc := documents.ClientListResponse{
		ItemsPerPage: 100,
		StartIndex:   1,
		TotalResults: len(cl),
		Schemas:      schemas,
		Resources:    []documents.ClientResponse{},
	}

	for _, client := range cl {
		doc.Resources = append(doc.Resources, client.ToDocument())
	}

	return doc
}
