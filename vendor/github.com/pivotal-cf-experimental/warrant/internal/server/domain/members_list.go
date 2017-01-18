package domain

import "github.com/pivotal-cf-experimental/warrant/internal/documents"

type MembersList []Member

func (ml MembersList) ToDocument() []documents.MemberResponse {
	members := []documents.MemberResponse{}

	for _, member := range ml {
		members = append(members, member.ToDocument())
	}

	return members
}
