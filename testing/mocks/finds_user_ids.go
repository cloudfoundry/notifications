package mocks

type FindsUserIDs struct {
	UserIDsBelongingToOrganizationCall struct {
		Receives struct {
			OrgGUID string
			Role    string
			Token   string
		}
		Returns struct {
			UserIDs []string
			Error   error
		}
	}

	UserIDsBelongingToScopeCall struct {
		Receives struct {
			Token string
			Scope string
		}
		Returns struct {
			UserIDs []string
			Error   error
		}
	}

	UserIDsBelongingToSpaceCall struct {
		Receives struct {
			SpaceGUID string
			Token     string
		}
		Returns struct {
			UserIDs []string
			Error   error
		}
	}
}

func NewFindsUserIDs() *FindsUserIDs {
	return &FindsUserIDs{}
}

func (f *FindsUserIDs) UserIDsBelongingToOrganization(orgGUID, role, token string) ([]string, error) {
	f.UserIDsBelongingToOrganizationCall.Receives.OrgGUID = orgGUID
	f.UserIDsBelongingToOrganizationCall.Receives.Role = role
	f.UserIDsBelongingToOrganizationCall.Receives.Token = token

	return f.UserIDsBelongingToOrganizationCall.Returns.UserIDs, f.UserIDsBelongingToOrganizationCall.Returns.Error
}

func (f *FindsUserIDs) UserIDsBelongingToScope(token, scope string) ([]string, error) {
	f.UserIDsBelongingToScopeCall.Receives.Token = token
	f.UserIDsBelongingToScopeCall.Receives.Scope = scope

	return f.UserIDsBelongingToScopeCall.Returns.UserIDs, f.UserIDsBelongingToScopeCall.Returns.Error
}

func (f *FindsUserIDs) UserIDsBelongingToSpace(spaceGUID, token string) ([]string, error) {
	f.UserIDsBelongingToSpaceCall.Receives.SpaceGUID = spaceGUID
	f.UserIDsBelongingToSpaceCall.Receives.Token = token

	return f.UserIDsBelongingToSpaceCall.Returns.UserIDs, f.UserIDsBelongingToSpaceCall.Returns.Error
}
