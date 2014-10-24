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

    Context("#Upsert", func() {
        Context("no error occurs", func() {
            var newTemplate models.Template
            var checkingTemplate models.Template

            BeforeEach(func() {
                newTemplate = models.Template{
                    Name:       "silly_template.user_body",
                    Text:       "omg",
                    HTML:       "<h1>OMG</h1>",
                    Overridden: true,
                }
            })

            It("inserts a template into the database", func() {
                var err error

                _, err = repo.Upsert(conn, newTemplate)
                count, err := conn.SelectInt("SELECT count(*) FROM `templates`")
                if err != nil {
                    panic(err)
                }

                Expect(count).To(Equal(int64(2)))
                err = conn.SelectOne(&checkingTemplate, "SELECT * FROM `templates` WHERE name=?", newTemplate.Name)
                if err != nil {
                    panic(err)
                }
                Expect(checkingTemplate.Text).To(Equal("omg"))
                Expect(checkingTemplate.HTML).To(Equal("<h1>OMG</h1>"))
            })

            It("updates a template currently in the database", func() {
                template.Text = "new text"
                _, err := repo.Upsert(conn, template)
                Expect(err).ToNot(HaveOccurred())

                count, err := conn.SelectInt("SELECT count(*) FROM `templates`")
                if err != nil {
                    panic(err)
                }
                Expect(count).To(Equal(int64(1)))

                err = conn.SelectOne(&checkingTemplate, "SELECT * FROM `templates` WHERE name=?", template.Name)
                if err != nil {
                    panic(err)
                }

                Expect(checkingTemplate.Text).To(Equal("new text"))
                Expect(checkingTemplate.HTML).To(Equal("<h1>containment unit breached!</h1>"))
            })
        })
    })
})
