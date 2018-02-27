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

type GroupsByDisplayName GroupsList

func (g GroupsByDisplayName) Len() int {
	return len(g)
}

func (g GroupsByDisplayName) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

func (g GroupsByDisplayName) Less(i, j int) bool {
	return g[i].DisplayName < g[j].DisplayName
}

type GroupsByCreatedAt GroupsList

func (g GroupsByCreatedAt) Len() int {
	return len(g)
}

func (g GroupsByCreatedAt) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

func (g GroupsByCreatedAt) Less(i, j int) bool {
	return g[i].CreatedAt.Before(g[j].CreatedAt)
}
