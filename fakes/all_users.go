package fakes

type AllUsers struct {
	AllUserGUIDsCall struct {
		Returns []string
		Error   error
	}
}

func NewAllUsers() *AllUsers {
	return &AllUsers{}
}

func (fake *AllUsers) AllUserGUIDs() ([]string, error) {
	return fake.AllUserGUIDsCall.Returns, fake.AllUserGUIDsCall.Error
}
