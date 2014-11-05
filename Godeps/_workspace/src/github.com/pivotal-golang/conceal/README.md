# Conceal

Conceal provides methods for encrypting and decrypting data. Its primary job is to encrypt and base64-encode a byte slice. It base64-decodes and decrypts data given a Cloak instantited with the same pin used to encrypt the data.

### Index
[func NewCloak(pin []byte) (Cloak, error)](#new-cloak)

[type Cloak struct](#cloak)

  - [func (cloak Cloak) Veil(data []byte) ([]byte, error)](#veil)

  - [func (cloak Cloak) Unveil(data []byte) ([]byte, error)](#unveil)

[type CipherLengthError struct{}](#cipher-length-error)

### Usage
<a name="new-cloak"></a>
####func NewCloak

```
func NewCloak(pin[]byte) (Cloak, error)
```
NewCloak returns a new Cloak with an AES cipher.Block. The pin is converted to a 16-byte slice and used to create the AES cipher.Block.

<a name="cloak"></a>
#### type Cloak

```
type Cloak struct {
    cipherBlock cipher.Block
}
```
A Cloak implements the CloakInterface by Veiling and Unveiling data.

<a name="veil"></a>
#### func Veil
```
func (cloak Cloak) Veil(data []byte) ([]byte, error)
```
Veil AES encrypts and base64 encodes data. The encryption is randomized, so the same data will not return the same encrypted slice of bytes each time.

<a name="unveil"></a>
#### func Unveil
```
func (cloak Cloak) Unveil(data []byte) ([]byte, error)
```
Unveil base64 decodes and AES decrypts data. Unveil will return an error if data is less than 16 bytes.

<a name="cipher-length-error"></a>
#### type CipherLengthError
```
type CipherLengthError struct{}
```
A CipherLengthError is returned if the data passed to Unveil is less than 16 bytes.