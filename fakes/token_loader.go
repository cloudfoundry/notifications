package fakes

type TokenLoader struct {
	Token         string
	LoadError     error
	LoadWasCalled bool
}

func NewTokenLoader() *TokenLoader {
	return &TokenLoader{}
}

func (fake *TokenLoader) Load() (string, error) {
	fake.LoadWasCalled = true
	return fake.Token, fake.LoadError
}
