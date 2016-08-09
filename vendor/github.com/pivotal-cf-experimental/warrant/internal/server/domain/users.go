package domain

type Users struct {
	store map[string]User
}

func NewUsers() *Users {
	return &Users{
		store: make(map[string]User),
	}
}

func (collection Users) Add(u User) {
	collection.store[u.ID] = u
}

func (collection Users) Update(u User) {
	collection.store[u.ID] = u
}

func (collection Users) Get(id string) (User, bool) {
	u, ok := collection.store[id]
	return u, ok
}

func (collection Users) GetByName(name string) (User, bool) {
	for _, u := range collection.store {
		if u.UserName == name {
			return u, true
		}
	}

	return User{}, false
}

func (collection Users) Delete(id string) bool {
	_, ok := collection.store[id]
	delete(collection.store, id)
	return ok
}

func (collection *Users) Clear() {
	collection.store = make(map[string]User)
}
