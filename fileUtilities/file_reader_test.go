package fileUtilities_test

import (
    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/fileUtilities"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("fileReader", func() {
    Describe("ReadFile", func() {
        It("returns a string of the file contents at the specified location", func() {
            env := config.NewEnvironment()
            path := env.RootPath + "/fileUtilities/fixtures/test.text"
            contents, err := fileUtilities.ReadFile(path)
            if err != nil {
                panic(err)
            }

            Expect(contents).To(Equal("We have some content\n\n\nAnd some more\n\n"))
        })
    })
})
