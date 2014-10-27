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
    var template models.Template

    BeforeEach(func() {
        TruncateTables()
        repo = models.NewTemplatesRepo()
        env := config.NewEnvironment()
        db := models.NewDatabase(env.DatabaseURL)
        conn = db.Connection()

        template = models.Template{
            Name:       "raptor_template",
            Text:       "run and hide",
            HTML:       "<h1>containment unit breached!</h1>",
            Overridden: true,
        }

        conn.Insert(&template)
    })

    Context("#Find", func() {
        Context("the template is in the database", func() {
            It("returns the template when it is found", func() {
                raptorTemplate, err := repo.Find(conn, "raptor_template")

                Expect(err).ToNot(HaveOccurred())
                Expect(raptorTemplate.Name).To(Equal("raptor_template"))
                Expect(raptorTemplate.Text).To(Equal("run and hide"))
                Expect(raptorTemplate.HTML).To(Equal("<h1>containment unit breached!</h1>"))
                Expect(raptorTemplate.Overridden).To(BeTrue())
            })
        })

        Context("the template is not in the database", func() {
            It("returns a record not found error", func() {
                sillyTemplate, err := repo.Find(conn, "silly_template")

                Expect(sillyTemplate).To(Equal(models.Template{}))
                Expect(err).To(HaveOccurred())
                Expect(err).To(MatchError(models.ErrRecordNotFound{}))
            })
        })
    })

    Describe("#Upsert", func() {
        It("inserts a template into the database", func() {
            newTemplate := models.Template{
                Name:       "silly_template.user_body",
                Text:       "omg",
                HTML:       "<h1>OMG</h1>",
                Overridden: true,
            }

            _, err := repo.Upsert(conn, newTemplate)
            Expect(err).ToNot(HaveOccurred())

            foundTemplate, err := repo.Find(conn, newTemplate.Name)
            if err != nil {
                panic(err)
            }

            Expect(foundTemplate.Text).To(Equal(newTemplate.Text))
            Expect(foundTemplate.HTML).To(Equal(newTemplate.HTML))
        })

        It("updates a template currently in the database", func() {
            template.Text = "new text"
            _, err := repo.Upsert(conn, template)
            Expect(err).ToNot(HaveOccurred())

            foundTemplate, err := repo.Find(conn, template.Name)
            if err != nil {
                panic(err)
            }

            Expect(foundTemplate.Text).To(Equal(template.Text))
            Expect(foundTemplate.HTML).To(Equal(template.HTML))
        })
    })

    Describe("#Destroy", func() {
        Context("the template exists in the database", func() {
            It("deletes the template", func() {
                _, err := repo.Find(conn, template.Name)
                if err != nil {
                    panic(err)
                }

                err = repo.Destroy(conn, template.Name)
                Expect(err).ToNot(HaveOccurred())

                _, err = repo.Find(conn, template.Name)
                Expect(err).To(Equal(models.ErrRecordNotFound{}))

            })
        })

        Context("the template does not exist in the database", func() {
            It("does not return an error", func() {
                err := repo.Destroy(conn, "knockknock")
                Expect(err).ToNot(HaveOccurred())

                _, err = repo.Find(conn, "knockknock")
                Expect(err).To(Equal(models.ErrRecordNotFound{}))

            })
        })
    })
})
