package fakes

type TokenLoader struct {
    Token     string
    LoadError error
}

func NewTokenLoader() *TokenLoader {
    return &TokenLoader{}
}

func (fake *TokenLoader) Load() (string, error) {
    return fake.Token, fake.LoadError
}
