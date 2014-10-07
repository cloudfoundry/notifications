package fakes

type FakeCryptoClient struct {
    EncryptedResult string
    EncryptArgument string
}

func (cryp *FakeCryptoClient) Encrypt(data string) (string, error) {
    cryp.EncryptArgument = data
    return cryp.EncryptedResult, nil
}

func (cryp *FakeCryptoClient) Decrypt(data string) (string, error) {
    return "stuff", nil
}
