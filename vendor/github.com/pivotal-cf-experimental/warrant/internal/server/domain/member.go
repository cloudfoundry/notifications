package domain

import (
	"github.com/pivotal-cf-experimental/warrant/internal/documents"
)

type Member struct {
	Origin string
	Type   string
	Value  string
}

func NewMemberFromDocument(document documents.CreateMemberRequest) Member {
	return Member{
		Type:   document.Type,
		Value:  document.Value,
		Origin: document.Origin,
	}
}

func (m Member) ToDocument() documents.MemberResponse {
	return documents.MemberResponse{
		Type:   m.Type,
		Value:  m.Value,
		Origin: m.Origin,
	}
}
