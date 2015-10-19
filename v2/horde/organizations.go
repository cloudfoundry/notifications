package horde

import (
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/pivotal-golang/lager"
)

type userFinder interface {
	UserIDsBelongingToOrganization(orgGUID, role, token string) (userGUIDs []string, err error)
	UserIDsBelongingToSpace(spaceGUID, token string) (userGUIDs []string, err error)
}

type orgFinder interface {
	Load(orgGUID, token string) (cf.CloudControllerOrganization, error)
}

type tokenLoader interface {
	Load(uaaHost string) (token string, err error)
}

type Organizations struct {
	userFinder  userFinder
	orgFinder   orgFinder
	tokenLoader tokenLoader
	uaaHost     string
}

func NewOrganizations(userFinder userFinder, orgFinder orgFinder, tokenLoader tokenLoader, uaaHost string) Organizations {
	return Organizations{
		userFinder:  userFinder,
		orgFinder:   orgFinder,
		tokenLoader: tokenLoader,
		uaaHost:     uaaHost,
	}
}

func (o Organizations) GenerateAudiences(orgGUIDs []string, logger lager.Logger) ([]Audience, error) {
	var audiences []Audience

	token, err := o.tokenLoader.Load(o.uaaHost)
	if err != nil {
		return audiences, err
	}

	for orgCounter, orgGUID := range orgGUIDs {
		var users []User

		if orgCounter%100 == 0 {
			logger.Debug("number of organizations", lager.Data{
				"processed": orgCounter,
			})
		}

		org, err := o.orgFinder.Load(orgGUID, token)
		if err != nil {
			if _, ok := err.(cf.NotFoundError); ok {
				continue
			}
			return audiences, err
		}

		userGUIDs, err := o.userFinder.UserIDsBelongingToOrganization(orgGUID, "", token)
		if err != nil {
			return audiences, err
		}

		for _, userGUID := range userGUIDs {
			users = append(users, User{GUID: userGUID})
		}

		audiences = append(audiences, Audience{
			Users:       users,
			Endorsement: fmt.Sprintf("You received this message because you belong to the %s organization.", org.Name),
		})
	}

	return audiences, nil
}
