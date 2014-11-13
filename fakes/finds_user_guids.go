package fakes

type FindsUserGUIDs struct {
	SpaceGuids                            map[string][]string
	UserGUIDsBelongingToSpaceError        error
	OrganizationGuids                     map[string][]string
	UserGUIDsBelongingToOrganizationError error
	GUIDsWithScopes                       map[string][]string
	UserGUIDsBelongingToScopeError        error
}

func NewFindsUserGUIDs() *FindsUserGUIDs {
	return &FindsUserGUIDs{
		SpaceGuids:        make(map[string][]string),
		OrganizationGuids: make(map[string][]string),
		GUIDsWithScopes:   make(map[string][]string),
	}
}

func (finder FindsUserGUIDs) UserGUIDsBelongingToSpace(spaceGUID, token string) ([]string, error) {
	return finder.SpaceGuids[spaceGUID], finder.UserGUIDsBelongingToSpaceError
}

func (finder FindsUserGUIDs) UserGUIDsBelongingToOrganization(orgGUID, role, token string) ([]string, error) {
	return finder.OrganizationGuids[orgGUID], finder.UserGUIDsBelongingToOrganizationError
}

func (finder FindsUserGUIDs) UserGUIDsBelongingToScope(scope string) ([]string, error) {
	return finder.GUIDsWithScopes[scope], finder.UserGUIDsBelongingToScopeError
}
