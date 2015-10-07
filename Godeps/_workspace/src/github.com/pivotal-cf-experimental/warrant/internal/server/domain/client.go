package domain

import (
	"errors"
	"fmt"

	"github.com/pivotal-cf-experimental/warrant/internal/documents"
)

var validGrantTypes = []string{"implicit", "refresh_token", "authorization_code", "client_credentials", "password"}

type client struct {
	ID                   string
	Secret               string
	Scope                []string
	ResourceIDs          []string
	Authorities          []string
	AuthorizedGrantTypes []string
	AccessTokenValidity  int
}

func NewClientFromDocument(document documents.CreateClientRequest) client {
	return client{
		ID:                   document.ClientID,
		Secret:               document.ClientSecret,
		Scope:                document.Scope,
		ResourceIDs:          document.ResourceIDs,
		Authorities:          document.Authorities,
		AuthorizedGrantTypes: document.AuthorizedGrantTypes,
		AccessTokenValidity:  document.AccessTokenValidity,
	}
}

func (c client) ToDocument() documents.ClientResponse {
	return documents.ClientResponse{
		ClientID:             c.ID,
		Scope:                shuffle(c.Scope),
		ResourceIDs:          shuffle(c.ResourceIDs),
		Authorities:          shuffle(c.Authorities),
		AuthorizedGrantTypes: shuffle(c.AuthorizedGrantTypes),
		AccessTokenValidity:  c.AccessTokenValidity,
	}
}

func (c client) Validate() error {
	for _, grantType := range c.AuthorizedGrantTypes {
		if !contains(validGrantTypes, grantType) {
			msg := fmt.Sprintf("%s is not an allowed grant type. Must be one of: %v", grantType, validGrantTypes)
			return errors.New(msg)
		}
	}

	return nil
}
