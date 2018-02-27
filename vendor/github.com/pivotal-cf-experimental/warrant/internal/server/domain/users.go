package domain

type Users struct {
	store map[string]User
}

func NewUsers() *Users {
	return &Users{
		store: make(map[string]User),
	}
}

func (collection Users) All() []User {
	var users []User
	for _, u := range collection.store {
		users = append(users, u)
	}
	return users
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

type ByEmail UsersList

func (users ByEmail) Len() int {
	return len(users)
}

func (users ByEmail) Swap(i, j int) {
	users[i], users[j] = users[j], users[i]
}

func (users ByEmail) Less(i, j int) bool {
	return users[i].Emails[0] < users[j].Emails[0]
}

type ByCreated UsersList

func (users ByCreated) Len() int {
	return len(users)
}

func (users ByCreated) Swap(i, j int) {
	users[i], users[j] = users[j], users[i]
}

func (users ByCreated) Less(i, j int) bool {
	return users[i].CreatedAt.Before(users[j].CreatedAt)
}
