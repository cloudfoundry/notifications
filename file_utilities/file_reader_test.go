package file_utilities_test

import (
    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/file_utilities"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("fileReader", func() {
    Describe("ReadFile", func() {
        It("returns a string of the file contents at the specified location", func() {
            env := config.NewEnvironment()
            path := env.RootPath + "/file_utilities/fixtures/test.text"
            contents, err := file_utilities.ReadFile(path)
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
            path = env.RootPath + "/file_utilities/fixtures/test.text"
        })

        It("returns true if the file exists", func() {
            response := file_utilities.FileExists(path)
            Expect(response).To(Equal(true))
        })

        It("returns false the file does not exist", func() {
            response := file_utilities.FileExists(path + "not.There")
            Expect(response).To(Equal(false))
        })

    })
})
