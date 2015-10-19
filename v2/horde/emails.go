package horde

import "github.com/pivotal-golang/lager"

type Emails struct {
}

func NewEmails() Emails {
	return Emails{}
}

func (e Emails) GenerateAudiences(emails []string, logger lager.Logger) ([]Audience, error) {
	var users []User
	for _, email := range emails {
		users = append(users, User{Email: email})
	}

	return []Audience{{
		Users:       users,
		Endorsement: "This message was sent directly to your email address.",
	}}, nil
}
