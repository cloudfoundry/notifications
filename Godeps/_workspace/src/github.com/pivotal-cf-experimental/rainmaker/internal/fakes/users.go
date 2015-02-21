package fakes

import "sort"

type Users struct {
	store map[string]User
}

func NewUsers() *Users {
	return &Users{
		store: make(map[string]User),
	}
}

func (users *Users) Add(usersToAdd ...User) {
	for _, user := range usersToAdd {
		users.store[user.GUID] = user
	}
}

func (users *Users) Get(guid string) (User, bool) {
	user, ok := users.store[guid]
	return user, ok
}

func (users *Users) Associate() {
}

func (users *Users) Clear() {
	users.store = make(map[string]User)
}

func (users Users) Len() int {
	return len(users.store)
}

func (users Users) Items() []interface{} {
	guids := sort.StringSlice([]string{})
	for _, user := range users.store {
		guids = append(guids, user.GUID)
	}

	sort.Sort(guids)

	var items []interface{}
	for _, guid := range guids {
		items = append(items, users.store[guid])
	}

	return items
}
