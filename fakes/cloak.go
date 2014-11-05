package fakes

type Cloak struct {
    EncryptedResult []byte
    DataToEncrypt   []byte
}

func (cloaker *Cloak) Veil(data []byte) ([]byte, error) {
    cloaker.DataToEncrypt = data
    return cloaker.EncryptedResult, nil
}

func (cloaker *Cloak) Unveil(data []byte) ([]byte, error) {
    return []byte("what"), nil
}
