package postal

type TokenLoader struct {
    uaaClient UAAInterface
}

func NewTokenLoader(uaaClient UAAInterface) TokenLoader {
    return TokenLoader{
        uaaClient: uaaClient,
    }
}

func (loader TokenLoader) Load() (string, error) {
    token, err := loader.uaaClient.GetClientToken()
    if err != nil {
        err = UAAErrorFor(err)
        return "", err
    }
    return token.Access, nil
}
