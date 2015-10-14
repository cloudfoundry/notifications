package horde

import (
	"fmt"

	"github.com/cloudfoundry-incubator/notifications/cf"
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
}

func NewSpaces(userFinder userFinder, orgFinder orgFinder, spaceFinder spaceFinder, tokenLoader tokenLoader, uaaHost string) Spaces {
	return Spaces{
		userFinder:  userFinder,
		orgFinder:   orgFinder,
		spaceFinder: spaceFinder,
		tokenLoader: tokenLoader,
		uaaHost:     uaaHost,
	}
}

func (s Spaces) GenerateAudiences(spaceGUIDs []string) ([]Audience, error) {
	var audiences []Audience

	token, err := s.tokenLoader.Load(s.uaaHost)
	if err != nil {
		return audiences, err
	}

	for _, spaceGUID := range spaceGUIDs {
		var users []User

		space, err := s.spaceFinder.Load(spaceGUID, token)
		if err != nil {
			return audiences, err
		}

		org, err := s.orgFinder.Load(space.OrganizationGUID, token)
		if err != nil {
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
