package models_test

import (
    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/models"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var _ = Describe("TemplatesRepo", func() {
    var repo models.TemplatesRepo
    var conn models.ConnectionInterface

    BeforeEach(func() {
        TruncateTables()
        repo = models.NewTemplatesRepo()
        env := config.NewEnvironment()
        db := models.NewDatabase(env.DatabaseURL)
        conn = db.Connection()

        template := &models.Template{
            Name:       "raptor_template",
            Text:       "run and hide",
            HTML:       "<h1>containment unit breached!</h1>",
            Overridden: true,
        }

        conn.Insert(template)
    })

    Context("#Find", func() {
        Context("no error occurs", func() {
            It("returns the template when it is found", func() {
                raptorTemplate, err := repo.Find(conn, "raptor_template")

                Expect(err).ToNot(HaveOccurred())
                Expect(raptorTemplate.Name).To(Equal("raptor_template"))
                Expect(raptorTemplate.Text).To(Equal("run and hide"))
                Expect(raptorTemplate.HTML).To(Equal("<h1>containment unit breached!</h1>"))
                Expect(raptorTemplate.Overridden).To(BeTrue())
            })
        })

        Context("an error occurs", func() {
            It("returns a record not found error and no template when the template is not in the DB", func() {
                sillyTemplate, err := repo.Find(conn, "silly_template")

                Expect(sillyTemplate).To(Equal(models.Template{}))
                Expect(err).To(HaveOccurred())
                Expect(err).To(MatchError(models.ErrRecordNotFound{}))
            })
        })
    })
})
