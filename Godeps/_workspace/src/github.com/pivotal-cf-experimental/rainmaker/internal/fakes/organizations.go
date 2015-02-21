package fakes

type Organizations struct {
	store map[string]Organization
}

func NewOrganizations() *Organizations {
	return &Organizations{
		store: make(map[string]Organization),
	}
}

func (orgs Organizations) Get(guid string) (Organization, bool) {
	org, ok := orgs.store[guid]
	return org, ok
}

func (orgs Organizations) Add(org Organization) {
	orgs.store[org.GUID] = org
}

func (orgs *Organizations) Clear() {
	orgs.store = make(map[string]Organization)
}
