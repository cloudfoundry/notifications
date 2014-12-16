package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TemplateAssigner", func() {
	Describe("AssignToClient", func() {
		var assigner services.TemplateAssigner
		var clientsRepo *fakes.ClientsRepo
		var templatesRepo *fakes.TemplatesRepo
		var conn *fakes.DBConn
		var database *fakes.Database

		BeforeEach(func() {
			var err error

			conn = fakes.NewDBConn()
			database = fakes.NewDatabase()

			clientsRepo = fakes.NewClientsRepo()
			_, err = clientsRepo.Create(conn, models.Client{
				ID: "my-client",
			})
			if err != nil {
				panic(err)
			}

			templatesRepo = fakes.NewTemplatesRepo()
			_, err = templatesRepo.Create(conn, models.Template{
				ID: "my-template",
			})
			if err != nil {
				panic(err)
			}

			assigner = services.NewTemplateAssigner(clientsRepo, templatesRepo, database)
		})

		It("assigns the template to the given client", func() {
			err := assigner.AssignToClient("my-client", "my-template")
			Expect(err).NotTo(HaveOccurred())

			client, err := clientsRepo.Find(conn, "my-client")
			if err != nil {
				panic(err)
			}

			Expect(client.Template).To(Equal("my-template"))
		})

		Context("when the request includes a non-existant id", func() {
			It("reports that the client cannot be found", func() {
				err := assigner.AssignToClient("bad-client", "my-template")
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(services.ClientMissingError("")))
			})

			It("reports that the template cannot be found", func() {
				err := assigner.AssignToClient("my-client", "non-existant-template")
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(services.TemplateAssignmentError("")))
			})
		})

		Context("when it gets an error it doesn't understand", func() {
			Context("on finding the client", func() {
				It("returns any errors it doesn't understand", func() {
					clientsRepo.FindError = errors.New("database connection failure")
					err := assigner.AssignToClient("my-client", "my-template")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("database connection failure"))
				})
			})
			Context("on finding the template", func() {
				It("returns any errors it doesn't understand (part 2)", func() {
					templatesRepo.FindError = errors.New("database failure")
					err := assigner.AssignToClient("my-client", "my-template")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("database failure"))

				})
			})

			Context("on updating the client", func() {
				It("Returns the error", func() {
					clientsRepo.UpdateError = errors.New("database fail")
					err := assigner.AssignToClient("my-client", "my-template")
					Expect(err).To(HaveOccurred())
				})
			})

		})
	})
})
