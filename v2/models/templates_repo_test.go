package models_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/testing/helpers"
	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v2/models"
	"github.com/nu7hatch/gouuid"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TemplatesRepo", func() {
	var (
		repo          models.TemplatesRepository
		conn          db.ConnectionInterface
		guidGenerator *mocks.GUIDGenerator
	)

	BeforeEach(func() {
		database := db.NewDatabase(sqlDB, db.Config{})
		helpers.TruncateTables(database)

		guid1 := uuid.UUID([16]byte{0xDE, 0xAD, 0xBE, 0xEF, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55})
		guid2 := uuid.UUID([16]byte{0xDE, 0xAD, 0xBE, 0xEF, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11, 0x22, 0x33, 0x44, 0x56})
		guidGenerator = mocks.NewGUIDGenerator()
		guidGenerator.GenerateCall.Returns.GUIDs = []*uuid.UUID{&guid1, &guid2}

		repo = models.NewTemplatesRepository(guidGenerator.Generate)
		conn = database.Connection()
	})

	Describe("Insert", func() {
		It("returns the data", func() {
			createdTemplate, err := repo.Insert(conn, models.Template{
				Name:     "some-template",
				ClientID: "some-client-id",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(createdTemplate.ID).To(Equal("deadbeef-aabb-ccdd-eeff-001122334455"))
		})

		Context("failure cases", func() {
			It("returns an error if it happens", func() {
				_, err := repo.Insert(conn, models.Template{
					Name:     "some-template",
					ClientID: "some-client-id",
				})
				Expect(err).NotTo(HaveOccurred())

				_, err = repo.Insert(conn, models.Template{
					Name:     "some-template",
					ClientID: "some-client-id",
				})
				Expect(err).To(BeAssignableToTypeOf(models.DuplicateRecordError{}))
			})
		})
	})

	Describe("Get", func() {
		It("fetches the template given a template_id", func() {
			createdTemplate, err := repo.Insert(conn, models.Template{
				Name:     "some-template",
				ClientID: "some-client-id",
			})
			Expect(err).NotTo(HaveOccurred())

			template, err := repo.Get(conn, createdTemplate.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(template).To(Equal(createdTemplate))
		})

		Context("failure cases", func() {
			It("returns not found error if it happens", func() {
				_, err := repo.Get(conn, "missing-template-id")
				Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError{}))
			})
		})
	})

	Describe("Delete", func() {
		It("deletes the template given a template_id", func() {
			template, err := repo.Insert(conn, models.Template{
				Name:     "some-template",
				ClientID: "some-client-id",
			})
			Expect(err).NotTo(HaveOccurred())

			err = repo.Delete(conn, template.ID)
			Expect(err).NotTo(HaveOccurred())

			_, err = repo.Get(conn, template.ID)
			Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError{}))
		})

		Context("failure cases", func() {
			It("returns not found error if it happens", func() {
				err := repo.Delete(conn, "missing-template-id")
				Expect(err).To(BeAssignableToTypeOf(models.RecordNotFoundError{}))
			})
		})
	})

	Describe("Update", func() {
		It("updates the template", func() {
			createdTemplate, err := repo.Insert(conn, models.Template{
				Name:     "some-template",
				ClientID: "some-client-id",
			})
			Expect(err).NotTo(HaveOccurred())

			createdTemplate.Name = "new-template"

			updatedTemplate, err := repo.Update(conn, createdTemplate)
			Expect(err).NotTo(HaveOccurred())
			Expect(updatedTemplate.ID).To(Equal(createdTemplate.ID))
			Expect(updatedTemplate.Name).To(Equal("new-template"))

			template, err := repo.Get(conn, updatedTemplate.ID)
			Expect(err).NotTo(HaveOccurred())
			Expect(template).To(Equal(updatedTemplate))
		})

		Context("failure cases", func() {
			It("returns an error when the update call fails", func() {
				fakeConn := mocks.NewConnection()
				fakeConn.UpdateCall.Returns.Error = errors.New("some db error")
				conn = fakeConn

				_, err := repo.Update(conn, models.Template{})
				Expect(err).To(MatchError(errors.New("some db error")))
			})
		})
	})
})
