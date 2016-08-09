package domain

type Clients struct {
	store map[string]client
}

func NewClients() *Clients {
	return &Clients{
		store: make(map[string]client),
	}
}

func (collection Clients) Add(c client) {
	collection.store[c.ID] = c
}

func (collection Clients) Get(id string) (client, bool) {
	c, ok := collection.store[id]
	return c, ok
}

func (collection *Clients) Clear() {
	collection.store = make(map[string]client)
}

func (collection Clients) Delete(id string) bool {
	_, ok := collection.store[id]
	delete(collection.store, id)
	return ok
}
