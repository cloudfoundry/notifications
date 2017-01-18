package domain

type Groups struct {
	store map[string]group
}

func NewGroups() *Groups {
	return &Groups{
		store: make(map[string]group),
	}
}

func (collection Groups) Add(g group) {
	collection.store[g.ID] = g
}

func (collection Groups) Get(id string) (group, bool) {
	g, ok := collection.store[id]
	return g, ok
}

func (collection Groups) Update(g group) {
	collection.store[g.ID] = g
}

func (collection Groups) All() []group {
	var groups []group
	for _, g := range collection.store {
		groups = append(groups, g)
	}

	return groups
}

func (collection Groups) Delete(id string) bool {
	_, ok := collection.store[id]
	delete(collection.store, id)
	return ok
}

func (collection *Groups) Clear() {
	collection.store = make(map[string]group)
}

func (collection Groups) AddMember(id string, member Member) (Member, bool) {
	group, ok := collection.store[id]
	if !ok {
		return Member{}, false
	}

	group.Members = append(group.Members, member)
	collection.store[id] = group

	return member, true
}

func (collection Groups) ListMembers(id string) ([]Member, bool) {
	group, ok := collection.store[id]
	if !ok {
		return []Member{}, false
	}

	var members []Member
	for _, m := range group.Members {
		members = append(members, m)
	}

	return members, true
}

func (collection Groups) CheckMembership(id, memberID string) (Member, bool) {
	g, ok := collection.store[id]
	if !ok {
		return Member{}, false
	}

	for _, member := range g.Members {
		if member.Value == memberID {
			return member, true
		}
	}

	return Member{}, false
}

func (collection Groups) GetByName(name string) (group, bool) {
	for _, g := range collection.store {
		if g.DisplayName == name {
			return g, true
		}
	}

	return group{}, false
}
