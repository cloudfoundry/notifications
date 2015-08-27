package mocks

type OrgFinder struct {
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

func NewOrgFinder() *OrgFinder {
	return &OrgFinder{}
}

func (f *OrgFinder) Exists(guid string) (bool, error) {
	f.ExistsCall.Receives.GUID = guid

	return f.ExistsCall.Returns.Exists, f.ExistsCall.Returns.Error
}
