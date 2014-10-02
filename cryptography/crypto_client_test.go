package cryptography_test

import (
    "regexp"

    "github.com/cloudfoundry-incubator/notifications/cryptography"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Cryptography", func() {

    Describe("NewURLCryptoClient", func() {
        It("returns an error when the length of the key is not 16 bytes", func() {
            _, err := cryptography.NewURLCryptoClient("tooShort")

            Expect(err).ToNot(BeNil())
            Expect(err).To(BeAssignableToTypeOf(cryptography.InvalidKeyError{}))
            Expect(err.Error()).To(Equal("The key needs to be 16 bytes"))
        })

        It("returns a CryptoClient when the key is valid", func() {
            client, err := cryptography.NewURLCryptoClient("sixteenbytes!!!!")
            if err != nil {
                panic(err)
            }

            Expect(client).To(BeAssignableToTypeOf(cryptography.URLCryptoClient{}))
        })
    })

    Describe("URLCryptoClient", func() {
        Describe("Encrypt And Decrypt", func() {
            It("2 way encrypts the string passed in", func() {
                client, err := cryptography.NewURLCryptoClient("sixteenbytes!!!!")
                if err != nil {
                    panic(err)
                }

                message := "this-message-is-secret"
                encryptedMessage, err := client.Encrypt(message)
                if err != nil {
                    panic(err)
                }

                decryptedMessage, err := client.Decrypt(encryptedMessage)
                if err != nil {
                    panic(err)
                }

                Expect(decryptedMessage).To(Equal(message))
            })

            It("The encryption only contains valid URL characters", func() {
                client, err := cryptography.NewURLCryptoClient("sixteenbytes!!!!")
                if err != nil {
                    panic(err)
                }

                message := "this-message-is-secret"
                encryptedMessage, err := client.Encrypt(message)
                if err != nil {
                    panic(err)
                }

                regex := regexp.MustCompile("^[a-zA-Z0-9-._~:/?#\\[\\]@!$&'()*+,;=]*$")
                results := regex.FindStringSubmatch(encryptedMessage)
                Expect(len(results)).To(Equal(1))
                Expect(results).ToNot(Equal(""))

            })

            Context("error cases", func() {
                It("the key and data are not the same byte length", func() {
                    client, err := cryptography.NewURLCryptoClient("sixteenbytes!!!!")
                    if err != nil {
                        panic(err)
                    }

                    _, err = client.Decrypt("oops")

                    Expect(err).ToNot(BeNil())
                    Expect(err).To(BeAssignableToTypeOf(cryptography.CipherLengthError{}))
                    Expect(err.Error()).To(Equal("Data length is too short"))
                })

                It("the data does not comply with base64 encoding", func() {
                    client, err := cryptography.NewURLCryptoClient("sixteenbytes!!!!")
                    if err != nil {
                        panic(err)
                    }

                    _, err = client.Decrypt("lsdkajfsaldkfjas;dlkfjasdlkfjdaslkjfadskljfasdlkjfas;dlkfjasldkjfasldkjfasld")

                    Expect(err).ToNot(BeNil())
                })
            })
        })
    })
})
