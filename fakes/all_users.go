package fakes

type AllUsers struct {
	LoadError error
	GUIDs     []string
}

func NewAllUsers() *AllUsers {
	return &AllUsers{}
}

func (fake *AllUsers) AllUserGUIDs() ([]string, error) {
	return fake.GUIDs, fake.LoadError
}
