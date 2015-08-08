package fakes

type AllUsers struct {
	Token            string
	AllUserGUIDsCall struct {
		Returns []string
		Error   error
	}
}

func NewAllUsers() *AllUsers {
	return &AllUsers{}
}

func (fake *AllUsers) AllUserGUIDs(token string) ([]string, error) {
	fake.Token = token
	return fake.AllUserGUIDsCall.Returns, fake.AllUserGUIDsCall.Error
}
