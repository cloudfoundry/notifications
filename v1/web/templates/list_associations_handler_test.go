package templates_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/collections"
	"github.com/cloudfoundry-incubator/notifications/v1/web/templates"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type TemplateAssociation struct {
	Client       string
	Notification string
}

var _ = Describe("ListAssociationsHandler", func() {
	var (
		handler     templates.ListAssociationsHandler
		templateID  string
		writer      *httptest.ResponseRecorder
		request     *http.Request
		lister      *mocks.TemplateAssociationLister
		errorWriter *mocks.ErrorWriter
		database    *mocks.Database
		connection  *mocks.Connection
		context     stack.Context
	)

	BeforeEach(func() {
		var err error

		templateID = "banana-template"
		lister = mocks.NewTemplateAssociationLister()
		lister.ListCall.Returns.Associations = []collections.TemplateAssociation{
			{
				ClientID: "some-client",
			},
			{
				ClientID:       "some-client",
				NotificationID: "some-notification",
			},
			{
				ClientID:       "another-client",
				NotificationID: "another-notification",
			},
		}

		errorWriter = mocks.NewErrorWriter()

		writer = httptest.NewRecorder()
		request, err = http.NewRequest("GET", "/templates/"+templateID+"/associations", nil)
		Expect(err).NotTo(HaveOccurred())

		connection = mocks.NewConnection()
		database = mocks.NewDatabase()
		database.ConnectionCall.Returns.Connection = connection

		context = stack.NewContext()
		context.Set("database", database)

		handler = templates.NewListAssociationsHandler(lister, errorWriter)
	})

	It("returns a list of clients and notifications associated to the given template", func() {
		handler.ServeHTTP(writer, request, context)

		Expect(writer.Code).To(Equal(http.StatusOK))

		var assoc struct {
			Associations []TemplateAssociation
		}
		err := json.Unmarshal(writer.Body.Bytes(), &assoc)
		Expect(err).NotTo(HaveOccurred())

		associations := assoc.Associations

		Expect(associations).To(HaveLen(3))
		Expect(associations).To(ContainElement(TemplateAssociation{
			Client: "some-client",
		}))
		Expect(associations).To(ContainElement(TemplateAssociation{
			Client:       "some-client",
			Notification: "some-notification",
		}))
		Expect(associations).To(ContainElement(TemplateAssociation{
			Client:       "another-client",
			Notification: "another-notification",
		}))

		Expect(lister.ListCall.Receives.Connection).To(Equal(connection))
		Expect(lister.ListCall.Receives.TemplateID).To(Equal(templateID))
	})

	Context("when errors occur", func() {
		Context("when the lister service returns an error", func() {
			It("delegates to the error handler", func() {
				lister.ListCall.Returns.Error = errors.New("db failed or something")

				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.WriteCall.Receives.Error).To(MatchError(errors.New("db failed or something")))
			})
		})
	})
})
