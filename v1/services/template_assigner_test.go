package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TemplateAssigner", func() {
	var assigner services.TemplateAssigner
	var kindsRepo *fakes.KindsRepo
	var clientsRepo *fakes.ClientsRepo
	var templatesRepo *fakes.TemplatesRepo
	var conn *fakes.Connection
	var database *fakes.Database

	BeforeEach(func() {
		conn = fakes.NewConnection()
		database = fakes.NewDatabase()
		clientsRepo = fakes.NewClientsRepo()
		kindsRepo = fakes.NewKindsRepo()
		templatesRepo = fakes.NewTemplatesRepo()
		assigner = services.NewTemplateAssigner(clientsRepo, kindsRepo, templatesRepo)
	})

	Describe("AssignToClient", func() {
		var client models.Client

		BeforeEach(func() {
			var err error

			client, err = clientsRepo.Create(conn, models.Client{
				ID: "my-client",
			})
			if err != nil {
				panic(err)
			}

			_, err = templatesRepo.Create(conn, models.Template{
				ID: "default",
			})
			if err != nil {
				panic(err)
			}

			_, err = templatesRepo.Create(conn, models.Template{
				ID: "my-template",
			})
			if err != nil {
				panic(err)
			}
		})

		It("assigns the template to the given client", func() {
			err := assigner.AssignToClient(database, "my-client", "my-template")
			Expect(err).NotTo(HaveOccurred())

			client, err := clientsRepo.Find(conn, "my-client")
			if err != nil {
				panic(err)
			}

			Expect(client.TemplateID).To(Equal("my-template"))
		})

		Context("when the request includes a non-existant id", func() {
			It("reports that the client cannot be found", func() {
				err := assigner.AssignToClient(database, "bad-client", "my-template")
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))
			})

			It("reports that the template cannot be found", func() {
				err := assigner.AssignToClient(database, "my-client", "non-existant-template")
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(services.TemplateAssignmentError("")))
			})
		})

		Context("when the request should reset the template assignment", func() {
			BeforeEach(func() {
				var err error
				client.TemplateID = "some-random-template"
				client, err = clientsRepo.Update(conn, client)
				if err != nil {
					panic(err)
				}
			})

			It("allows template id of empty string to reset the assignment", func() {
				err := assigner.AssignToClient(database, "my-client", "")
				Expect(err).NotTo(HaveOccurred())

				client, err = clientsRepo.Find(conn, "my-client")
				if err != nil {
					panic(err)
				}

				Expect(client.TemplateID).To(Equal(models.DefaultTemplateID))
			})

			It("allows template id of default template id to reset the assignment", func() {
				err := assigner.AssignToClient(database, "my-client", models.DefaultTemplateID)
				Expect(err).NotTo(HaveOccurred())

				client, err = clientsRepo.Find(conn, "my-client")
				if err != nil {
					panic(err)
				}

				Expect(client.TemplateID).To(Equal(models.DefaultTemplateID))
			})
		})

		Context("when it gets an error it doesn't understand", func() {
			Context("on finding the client", func() {
				It("returns any errors it doesn't understand", func() {
					clientsRepo.FindCall.Error = errors.New("database connection failure")
					err := assigner.AssignToClient(database, "my-client", "my-template")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("database connection failure"))
				})
			})
			Context("on finding the template", func() {
				It("returns any errors it doesn't understand (part 2)", func() {
					templatesRepo.FindError = errors.New("database failure")
					err := assigner.AssignToClient(database, "my-client", "my-template")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("database failure"))

				})
			})

			Context("on updating the client", func() {
				It("Returns the error", func() {
					clientsRepo.UpdateCall.Error = errors.New("database fail")
					err := assigner.AssignToClient(database, "my-client", "my-template")
					Expect(err).To(HaveOccurred())
				})
			})

		})
	})

	Describe("AssignToNotification", func() {
		var kind models.Kind

		BeforeEach(func() {
			client, err := clientsRepo.Create(conn, models.Client{
				ID: "my-client",
			})
			if err != nil {
				panic(err)
			}

			_, err = templatesRepo.Create(conn, models.Template{
				ID: "default",
			})
			if err != nil {
				panic(err)
			}

			kind, err = kindsRepo.Create(conn, models.Kind{
				ID:       "my-kind",
				ClientID: client.ID,
			})
			if err != nil {
				panic(err)
			}

			_, err = templatesRepo.Create(conn, models.Template{
				ID: "my-template",
			})
			if err != nil {
				panic(err)
			}
		})

		It("assigns the template to the given kind", func() {
			err := assigner.AssignToNotification(database, "my-client", "my-kind", "my-template")
			Expect(err).NotTo(HaveOccurred())

			kind, err = kindsRepo.Find(conn, "my-kind", "my-client")
			if err != nil {
				panic(err)
			}

			Expect(kind.TemplateID).To(Equal("my-template"))
		})

		Context("when the request includes a non-existant id", func() {
			It("reports that the client cannot be found", func() {
				err := assigner.AssignToNotification(database, "bad-client", "my-kind", "my-template")
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))
			})

			It("reports that the kind cannot be found", func() {
				err := assigner.AssignToNotification(database, "my-client", "bad-kind", "my-template")
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError("")))
			})

			It("reports that the template cannot be found", func() {
				err := assigner.AssignToNotification(database, "my-client", "my-kind", "non-existant-template")
				Expect(err).To(HaveOccurred())
				Expect(err).To(BeAssignableToTypeOf(services.TemplateAssignmentError("")))
			})
		})

		Context("when the request should reset the template assignment", func() {
			BeforeEach(func() {
				var err error
				kind.TemplateID = "some-random-template"
				kind, err = kindsRepo.Update(conn, kind)
				if err != nil {
					panic(err)
				}
			})

			It("allows template id of empty string to reset the assignment", func() {
				err := assigner.AssignToNotification(database, "my-client", "my-kind", "")
				Expect(err).NotTo(HaveOccurred())

				kind, err = kindsRepo.Find(conn, "my-kind", "my-client")
				if err != nil {
					panic(err)
				}

				Expect(kind.TemplateID).To(Equal(models.DefaultTemplateID))
			})

			It("allows template id of default template id to reset the assignment", func() {
				err := assigner.AssignToNotification(database, "my-client", "my-kind", models.DefaultTemplateID)
				Expect(err).NotTo(HaveOccurred())

				kind, err = kindsRepo.Find(conn, "my-kind", "my-client")
				if err != nil {
					panic(err)
				}

				Expect(kind.TemplateID).To(Equal(models.DefaultTemplateID))
			})
		})

		Context("when it gets an error it doesn't understand", func() {
			Context("on finding the client", func() {
				It("returns any errors it doesn't understand", func() {
					clientsRepo.FindCall.Error = errors.New("database connection failure")
					err := assigner.AssignToNotification(database, "my-client", "my-kind", "my-template")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("database connection failure"))
				})
			})
			Context("on finding the template", func() {
				It("returns any errors it doesn't understand (part 2)", func() {
					templatesRepo.FindError = errors.New("database failure")
					err := assigner.AssignToNotification(database, "my-client", "my-kind", "my-template")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("database failure"))

				})
			})

			Context("on updating the client", func() {
				It("Returns the error", func() {
					kindsRepo.UpdateError = errors.New("database fail")
					err := assigner.AssignToNotification(database, "my-client", "my-kind", "my-template")
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})
