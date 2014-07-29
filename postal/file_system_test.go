package postal_test

import (
    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/postal"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

const FixtureFile = "/postal/fixtures/test.text"

var _ = Describe("FileSystem", func() {
    var fs postal.FileSystem

    Describe("Read", func() {
        It("returns a string of the file contents at the specified location", func() {
            env := config.NewEnvironment()
            path := env.RootPath + FixtureFile
            contents, err := fs.Read(path)
            if err != nil {
                panic(err)
            }

            Expect(contents).To(Equal("We have some content\n\n\nAnd some more\n\n"))
        })
    })

    Describe("FileExists", func() {
        var path string

        BeforeEach(func() {
            env := config.NewEnvironment()
            path = env.RootPath + FixtureFile
        })

        It("returns true if the file exists", func() {
            response := fs.Exists(path)
            Expect(response).To(BeTrue())
        })

        It("returns false if the file does not exist", func() {
            response := fs.Exists(path + "not.There")
            Expect(response).To(BeFalse())
        })
    })
})
