package domain

import "github.com/pivotal-cf-experimental/warrant/internal/documents"

type GroupsList []group

func (gl GroupsList) ToDocument() documents.GroupListResponse {
	doc := documents.GroupListResponse{
		ItemsPerPage: 100,
		StartIndex:   1,
		TotalResults: len(gl),
		Schemas:      schemas,
	}

	for _, group := range gl {
		doc.Resources = append(doc.Resources, group.ToDocument())
	}

	return doc
}
