package mocks

type Cloak struct {
	VeilCall struct {
		Receives struct {
			PlainText []byte
		}
		Returns struct {
			CipherText []byte
			Error      error
		}
	}

	UnveilCall struct {
		Receives struct {
			CipherText []byte
		}
		Returns struct {
			PlainText []byte
			Error     error
		}
	}
}

func NewCloak() *Cloak {
	return &Cloak{}
}

func (c *Cloak) Veil(plaintext []byte) ([]byte, error) {
	c.VeilCall.Receives.PlainText = plaintext

	return c.VeilCall.Returns.CipherText, c.VeilCall.Returns.Error
}

func (c *Cloak) Unveil(ciphertext []byte) ([]byte, error) {
	c.UnveilCall.Receives.CipherText = ciphertext

	return c.UnveilCall.Returns.PlainText, c.UnveilCall.Returns.Error
}
