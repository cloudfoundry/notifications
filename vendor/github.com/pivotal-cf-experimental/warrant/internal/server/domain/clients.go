package domain

var ADMIN_CLIENT = Client{
	ID:     "admin",
	Name:   "admin",
	Secret: "admin",
	Scope:  []string{},
	ResourceIDs: []string{
		"clients", //TODO: This is not needed, but checked in some handlers incorrectly
		"password",
		"scim",
	},
	Authorities: []string{
		"clients.read",
		"clients.write",
		"clients.secret",
		"password.write",
		"uaa.admin",
		"scim.read",
		"scim.write",
	},
	AuthorizedGrantTypes: []string{
		"client_credentials", //TODO: we aren't checking that the client has this value when we generate tokens in the handlers
	},
	AccessTokenValidity: 3600,
	RedirectURI:         []string{},
	Autoapprove:         []string{},
}

type Clients struct {
	store map[string]Client
}

func NewClients() *Clients {
	return &Clients{
		store: map[string]Client{
			"admin": ADMIN_CLIENT,
		},
	}
}

func (collection Clients) All() []Client {
	var clients []Client
	for _, c := range collection.store {
		clients = append(clients, c)
	}
	return clients
}

func (collection Clients) Add(c Client) {
	collection.store[c.ID] = c
}

func (collection Clients) Get(id string) (Client, bool) {
	c, ok := collection.store[id]
	return c, ok
}

func (collection *Clients) Clear() {
	collection.store = map[string]Client{
		"admin": ADMIN_CLIENT,
	}
}

func (collection Clients) Delete(id string) bool {
	_, ok := collection.store[id]
	delete(collection.store, id)
	return ok
}

type ByName ClientsList

func (clients ByName) Len() int {
	return len(clients)
}

func (clients ByName) Swap(i, j int) {
	clients[i], clients[j] = clients[j], clients[i]
}

func (clients ByName) Less(i, j int) bool {
	return clients[i].Name < clients[j].Name
}

type ByID ClientsList

func (clients ByID) Len() int {
	return len(clients)
}

func (clients ByID) Swap(i, j int) {
	clients[i], clients[j] = clients[j], clients[i]
}

func (clients ByID) Less(i, j int) bool {
	return clients[i].ID < clients[j].ID
}
