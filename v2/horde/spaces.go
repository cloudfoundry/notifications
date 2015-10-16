package horde

import (
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/cf"
	"github.com/pivotal-golang/lager"
)

type spaceFinder interface {
	Load(spaceGUID, token string) (cf.CloudControllerSpace, error)
}

type Spaces struct {
	userFinder  userFinder
	orgFinder   orgFinder
	spaceFinder spaceFinder
	tokenLoader tokenLoader
	uaaHost     string
	logger      lager.Logger
}

func NewSpaces(userFinder userFinder, orgFinder orgFinder, spaceFinder spaceFinder, tokenLoader tokenLoader, uaaHost string, logger lager.Logger) Spaces {
	return Spaces{
		userFinder:  userFinder,
		orgFinder:   orgFinder,
		spaceFinder: spaceFinder,
		tokenLoader: tokenLoader,
		uaaHost:     uaaHost,
		logger:      logger,
	}
}

func (s Spaces) GenerateAudiences(spaceGUIDs []string) ([]Audience, error) {
	var audiences []Audience

	token, err := s.tokenLoader.Load(s.uaaHost)
	if err != nil {
		return audiences, err
	}

	var spaceCounter int
	for _, spaceGUID := range spaceGUIDs {
		var users []User

		if spaceCounter%100 == 0 {
			s.logger.Debug("number of spaces", lager.Data{
				"processed": spaceCounter,
			})
		}

		spaceCounter++

		space, err := s.spaceFinder.Load(spaceGUID, token)
		if err != nil {
			if _, ok := err.(cf.NotFoundError); ok {
				continue
			}
			return audiences, err
		}

		org, err := s.orgFinder.Load(space.OrganizationGUID, token)
		if err != nil {
			if _, ok := err.(cf.NotFoundError); ok {
				continue
			}
			return audiences, err
		}

		userGUIDs, err := s.userFinder.UserIDsBelongingToSpace(space.GUID, token)
		if err != nil {
			return audiences, err
		}

		for _, userGUID := range userGUIDs {
			users = append(users, User{GUID: userGUID})
		}

		audiences = append(audiences, Audience{
			Users:       users,
			Endorsement: fmt.Sprintf("You received this message because you belong to the %q space in the %q organization.", space.Name, org.Name),
		})
	}

	return audiences, nil
}
