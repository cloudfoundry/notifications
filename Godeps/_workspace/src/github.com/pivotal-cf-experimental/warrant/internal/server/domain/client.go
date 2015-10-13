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
	RedirectURI          []string
	Autoapprove          []string
}

func NewClientFromDocument(document documents.CreateUpdateClientRequest) client {
	return client{
		ID:                   document.ClientID,
		Secret:               document.ClientSecret,
		Scope:                document.Scope,
		ResourceIDs:          document.ResourceIDs,
		Authorities:          document.Authorities,
		AuthorizedGrantTypes: document.AuthorizedGrantTypes,
		AccessTokenValidity:  document.AccessTokenValidity,
		RedirectURI:          document.RedirectURI,
		Autoapprove:          document.Autoapprove,
	}
}

func (c client) ToDocument() documents.ClientResponse {
	return documents.ClientResponse{
		ClientID:             c.ID,
		Scope:                shuffle(c.Scope),
		ResourceIDs:          shuffle(c.ResourceIDs),
		Authorities:          shuffle(c.Authorities),
		AuthorizedGrantTypes: shuffle(c.AuthorizedGrantTypes),
		Autoapprove:          shuffle(c.Autoapprove),
		AccessTokenValidity:  c.AccessTokenValidity,
		RedirectURI:          c.RedirectURI,
	}
}

func (c client) Validate() error {
	for _, grantType := range c.AuthorizedGrantTypes {
		if !contains(validGrantTypes, grantType) {
			msg := fmt.Sprintf("%s is not an allowed grant type. Must be one of: %v", grantType, validGrantTypes)
			return errors.New(msg)
		}
	}

	if len(c.RedirectURI) > 0 {
		if !contains(c.AuthorizedGrantTypes, "implicit") && !contains(c.AuthorizedGrantTypes, "authorization_code") {
			msg := "A redirect_uri can only be used by implicit or authorization_code grant types."
			return errors.New(msg)
		}
	}

	if contains(c.AuthorizedGrantTypes, "implicit") && c.Secret != "" {
		msg := "Implicit grant should not have a client_secret"
		return errors.New(msg)
	}

	return nil
}
