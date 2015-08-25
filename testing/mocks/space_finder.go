package mocks

type SpaceFinder struct {
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

func NewSpaceFinder() *SpaceFinder {
	return &SpaceFinder{}
}

func (u *SpaceFinder) Exists(guid string) (bool, error) {
	u.ExistsCall.Receives.GUID = guid

	return u.ExistsCall.Returns.Exists, u.ExistsCall.Returns.Error
}
