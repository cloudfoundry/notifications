package mocks

type FindsUserGUIDs struct {
	UserGUIDsBelongingToOrganizationCall struct {
		Receives struct {
			OrgGUID string
			Role    string
			Token   string
		}
		Returns struct {
			UserGUIDs []string
			Error     error
		}
	}

	UserGUIDsBelongingToScopeCall struct {
		Receives struct {
			Token string
			Scope string
		}
		Returns struct {
			UserGUIDs []string
			Error     error
		}
	}

	UserGUIDsBelongingToSpaceCall struct {
		Receives struct {
			SpaceGUID string
			Token     string
		}
		Returns struct {
			UserGUIDs []string
			Error     error
		}
	}
}

func NewFindsUserGUIDs() *FindsUserGUIDs {
	return &FindsUserGUIDs{}
}

func (f *FindsUserGUIDs) UserGUIDsBelongingToOrganization(orgGUID, role, token string) ([]string, error) {
	f.UserGUIDsBelongingToOrganizationCall.Receives.OrgGUID = orgGUID
	f.UserGUIDsBelongingToOrganizationCall.Receives.Role = role
	f.UserGUIDsBelongingToOrganizationCall.Receives.Token = token

	return f.UserGUIDsBelongingToOrganizationCall.Returns.UserGUIDs, f.UserGUIDsBelongingToOrganizationCall.Returns.Error
}

func (f *FindsUserGUIDs) UserGUIDsBelongingToScope(token, scope string) ([]string, error) {
	f.UserGUIDsBelongingToScopeCall.Receives.Token = token
	f.UserGUIDsBelongingToScopeCall.Receives.Scope = scope

	return f.UserGUIDsBelongingToScopeCall.Returns.UserGUIDs, f.UserGUIDsBelongingToScopeCall.Returns.Error
}

func (f *FindsUserGUIDs) UserGUIDsBelongingToSpace(spaceGUID, token string) ([]string, error) {
	f.UserGUIDsBelongingToSpaceCall.Receives.SpaceGUID = spaceGUID
	f.UserGUIDsBelongingToSpaceCall.Receives.Token = token

	return f.UserGUIDsBelongingToSpaceCall.Returns.UserGUIDs, f.UserGUIDsBelongingToSpaceCall.Returns.Error
}
