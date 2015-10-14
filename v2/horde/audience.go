package horde

type Audience struct {
	Users       []User
	Endorsement string
}

type User struct {
	Email string
	GUID  string
}
