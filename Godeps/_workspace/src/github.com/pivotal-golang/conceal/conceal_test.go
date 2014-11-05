package conceal_test

import (
    "regexp"

    "github.com/pivotal-golang/conceal"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("Cloak", func() {

    Describe("NewCloak", func() {
        It("returns a Cloak when the key is valid", func() {
            cloak, err := conceal.NewCloak([]byte("sixteenbytes!!!!"))
            if err != nil {
                panic(err)
            }

            Expect(cloak).To(BeAssignableToTypeOf(conceal.Cloak{}))
        })
    })

    Describe("Veil And Unveil", func() {
        It("2 way encrypts the string passed in", func() {
            cloak, err := conceal.NewCloak([]byte("sixteenbytes!!!!"))
            if err != nil {
                panic(err)
            }

            message := []byte("this-message-is-secret")
            encryptedMessage, err := cloak.Veil([]byte(message))
            if err != nil {
                panic(err)
            }

            decryptedMessage, err := cloak.Unveil([]byte(encryptedMessage))
            if err != nil {
                panic(err)
            }

            Expect(decryptedMessage).To(Equal(message))
        })

        It("The encryption only contains valid characters", func() {
            cloak, err := conceal.NewCloak([]byte("sixteenbytes!!!!"))
            if err != nil {
                panic(err)
            }

            message := "this-message-is-secret"
            encryptedMessage, err := cloak.Veil([]byte(message))
            if err != nil {
                panic(err)
            }

            regex := regexp.MustCompile("^[a-zA-Z0-9-._~:/?#\\[\\]@!$&'()*+,;=]*$")
            results := regex.FindStringSubmatch(string(encryptedMessage))
            Expect(len(results)).To(Equal(1))
            Expect(results).ToNot(Equal(""))

        })

        It("Veil and unveil still works with a small key", func() {
            cloak, err := conceal.NewCloak([]byte("small"))
            if err != nil {
                panic(err)
            }

            message := []byte("this-message-is-secret")
            encryptedMessage, err := cloak.Veil([]byte(message))
            if err != nil {
                panic(err)
            }

            decryptedMessage, err := cloak.Unveil([]byte(encryptedMessage))
            if err != nil {
                panic(err)
            }

            Expect(decryptedMessage).To(Equal(message))
        })

        Context("error cases", func() {
            It("the data is not long enough", func() {
                cloak, err := conceal.NewCloak([]byte("sixteenbytes!!!!"))
                if err != nil {
                    panic(err)
                }

                _, err = cloak.Unveil([]byte("oops"))

                Expect(err).To(BeAssignableToTypeOf(conceal.CipherLengthError{}))
                Expect(err.Error()).To(Equal("Data length should be at least 16 bytes"))
            })

            It("the data does not comply with base64 encoding", func() {
                cloak, err := conceal.NewCloak([]byte("sixteenbytes!!!!"))
                if err != nil {
                    panic(err)
                }

                _, err = cloak.Unveil([]byte("lsdkajfsaldkfjas;dlkfjasdlkfjdaslkjfadskljfasdlkjfas;dlkfjasldkjfasldkjfasld"))

                Expect(err).ToNot(BeNil())
            })
        })
    })
})
