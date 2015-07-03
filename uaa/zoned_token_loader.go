package uaa

type ZonedTokenLoader struct {
	uaaClient UAAClientInterface
}

type UAAClientInterface interface {
	ZonedGetClientToken(string) (string, error)
}

func NewZonedTokenLoader(uaaClient UAAClientInterface) *ZonedTokenLoader {
	return &ZonedTokenLoader{
		uaaClient: uaaClient,
	}
}

func (z *ZonedTokenLoader) Load(uaaHost string) (string, error) {
	token, err := z.uaaClient.ZonedGetClientToken(uaaHost)
	return token, err
}
