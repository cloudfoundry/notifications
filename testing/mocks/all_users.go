package mocks

type AllUsers struct {
	AllUserGUIDsCall struct {
		Receives struct {
			Token string
		}
		Returns struct {
			GUIDs []string
			Error error
		}
	}
}

func NewAllUsers() *AllUsers {
	return &AllUsers{}
}

func (au *AllUsers) AllUserGUIDs(token string) ([]string, error) {
	au.AllUserGUIDsCall.Receives.Token = token
	return au.AllUserGUIDsCall.Returns.GUIDs, au.AllUserGUIDsCall.Returns.Error
}
