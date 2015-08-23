package mocks

type TokenLoader struct {
	LoadCall struct {
		Receives struct {
			UAAHost string
		}
		Returns struct {
			Token string
			Error error
		}
	}
}

func NewTokenLoader() *TokenLoader {
	return &TokenLoader{}
}

func (t *TokenLoader) Load(uaaHost string) (string, error) {
	t.LoadCall.Receives.UAAHost = uaaHost

	return t.LoadCall.Returns.Token, t.LoadCall.Returns.Error
}
