package fakes

type Spaces struct {
	store map[string]Space
}

func NewSpaces() *Spaces {
	return &Spaces{
		store: make(map[string]Space),
	}
}

func (spaces Spaces) Get(guid string) (Space, bool) {
	space, ok := spaces.store[guid]
	return space, ok
}

func (spaces Spaces) Add(space Space) {
	spaces.store[space.GUID] = space
}

func (spaces *Spaces) Clear() {
	spaces.store = make(map[string]Space)
}
