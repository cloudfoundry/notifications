package fakes

type FindsUserGUIDs struct {
    SpaceGuids                            map[string][]string
    UserGUIDsBelongingToSpaceError        error
    OrganizationGuids                     map[string][]string
    UserGUIDsBelongingToOrganizationError error
}

func NewFindsUserGUIDs() *FindsUserGUIDs {
    return &FindsUserGUIDs{
        SpaceGuids:        make(map[string][]string),
        OrganizationGuids: make(map[string][]string),
    }
}

func (finder FindsUserGUIDs) UserGUIDsBelongingToSpace(spaceGUID, token string) ([]string, error) {
    return finder.SpaceGuids[spaceGUID], finder.UserGUIDsBelongingToSpaceError
}

func (finder FindsUserGUIDs) UserGUIDsBelongingToOrganization(orgGUID, role, token string) ([]string, error) {
    return finder.OrganizationGuids[orgGUID], finder.UserGUIDsBelongingToOrganizationError
}
