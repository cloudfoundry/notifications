package domain

import (
	"errors"
	"fmt"

	"github.com/pivotal-cf-experimental/warrant/internal/documents"
)

var validGrantTypes = []string{
	"implicit",
	"refresh_token",
	"authorization_code",
	"client_credentials",
	"password",
}

type Client struct {
	ID                   string
	Name                 string
	Secret               string
	Scope                []string
	ResourceIDs          []string
	Authorities          []string
	AuthorizedGrantTypes []string
	AccessTokenValidity  int
	RedirectURI          []string
	Autoapprove          []string
}

func NewClientFromDocument(document documents.CreateUpdateClientRequest) Client {
	return Client{
		ID:                   document.ClientID,
		Name:                 document.Name,
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

func (c Client) ToDocument() documents.ClientResponse {
	return documents.ClientResponse{
		ClientID:             c.ID,
		Name:                 c.Name,
		Scope:                shuffle(c.Scope),
		ResourceIDs:          c.ResourceIDs,
		Authorities:          c.Authorities,
		AuthorizedGrantTypes: shuffle(c.AuthorizedGrantTypes),
		Autoapprove:          shuffle(c.Autoapprove),
		AccessTokenValidity:  c.AccessTokenValidity,
		RedirectURI:          c.RedirectURI,
	}
}

func (c Client) Validate() error {
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
