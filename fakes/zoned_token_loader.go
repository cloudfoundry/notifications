package fakes

type ZonedTokenLoader struct {
	Token        string
	LoadArgument string
	LoadError    error
}

func NewZonedTokenLoader() *ZonedTokenLoader {
	return &ZonedTokenLoader{}
}

func (z *ZonedTokenLoader) Load(uaaHost string) (string, error) {
	z.LoadArgument = uaaHost
	return z.Token, z.LoadError
}
