package mocks

type UserFinder struct {
	ExistsCall struct {
		Receives struct {
			GUID string
		}

		Returns struct {
			Exists bool
			Error  error
		}
	}
}

func NewUserFinder() *UserFinder {
	return &UserFinder{}
}

func (u *UserFinder) Exists(guid string) (bool, error) {
	u.ExistsCall.Receives.GUID = guid

	return u.ExistsCall.Returns.Exists, u.ExistsCall.Returns.Error
}
