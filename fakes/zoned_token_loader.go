package fakes

type TokenLoader struct {
	Token        string
	LoadArgument string
	LoadError    error
}

func NewTokenLoader() *TokenLoader {
	return &TokenLoader{}
}

func (t *TokenLoader) Load(uaaHost string) (string, error) {
	t.LoadArgument = uaaHost
	return t.Token, t.LoadError
}
