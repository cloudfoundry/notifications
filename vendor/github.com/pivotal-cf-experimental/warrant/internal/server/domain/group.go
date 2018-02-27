package domain

import (
	"time"

	"github.com/pivotal-cf-experimental/warrant/internal/documents"
	"github.com/pivotal-cf-experimental/warrant/internal/server/common"
)

type group struct {
	ID          string
	DisplayName string
	Description string
	Members     []Member
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Version     int
}

func NewGroupFromCreateDocument(request documents.CreateGroupRequest) group {
	var members []Member
	for _, member := range request.Members {
		members = append(members, Member{
			Value:  member.Value,
			Type:   member.Type,
			Origin: member.Origin,
		})
	}

	now := time.Now().UTC()
	id, err := common.NewUUID()
	if err != nil {
		panic(err)
	}

	return group{
		ID:          id,
		DisplayName: request.DisplayName,
		Description: request.Description,
		Members:     members,
		CreatedAt:   now,
		UpdatedAt:   now,
		Version:     0,
	}
}

func NewGroupFromUpdateDocument(request documents.CreateUpdateGroupRequest) group {
	var members []Member
	for _, member := range request.Members {
		members = append(members, Member{
			Value:  member.Value,
			Type:   member.Type,
			Origin: member.Origin,
		})
	}

	return group{
		ID:          request.ID,
		DisplayName: request.DisplayName,
		Description: request.Description,
		Members:     members,
		CreatedAt:   request.Meta.Created,
		UpdatedAt:   request.Meta.LastModified,
		Version:     request.Meta.Version,
	}
}

func (g group) ToDocument() documents.GroupResponse {
	var members []documents.MemberResponse
	for _, member := range g.Members {
		members = append(members, documents.MemberResponse{
			Value:  member.Value,
			Type:   member.Type,
			Origin: member.Origin,
		})
	}

	return documents.GroupResponse{
		Schemas:     schemas,
		ID:          g.ID,
		Description: g.Description,
		DisplayName: g.DisplayName,
		Members:     members,
		Meta: documents.Meta{
			Version:      g.Version,
			Created:      g.CreatedAt,
			LastModified: g.UpdatedAt,
		},
	}
}
