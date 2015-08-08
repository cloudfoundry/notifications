package templates_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/testing/fakes"
	"github.com/cloudfoundry-incubator/notifications/v1/services"
	"github.com/cloudfoundry-incubator/notifications/v1/web/templates"
	"github.com/ryanmoran/stack"

	. "github.com/onsi/ginkgo"
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
		lister      *fakes.TemplateAssociationLister
		errorWriter *fakes.ErrorWriter
		database    *fakes.Database
		context     stack.Context
	)

	BeforeEach(func() {
		var err error

		templateID = "banana-template"
		lister = fakes.NewTemplateAssociationLister()
		lister.Associations[templateID] = []services.TemplateAssociation{
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

		errorWriter = fakes.NewErrorWriter()

		writer = httptest.NewRecorder()
		request, err = http.NewRequest("GET", "/templates/"+templateID+"/associations", nil)
		Expect(err).NotTo(HaveOccurred())

		database = fakes.NewDatabase()
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

		Expect(lister.ListCall.Arguments).To(ConsistOf([]interface{}{database, templateID}))
	})

	Context("when errors occur", func() {
		Context("when the lister service returns an error", func() {
			It("delegates to the error handler", func() {
				lister.ListCall.Error = errors.New("db failed or something")

				handler.ServeHTTP(writer, request, context)
				Expect(errorWriter.WriteCall.Receives.Error).To(MatchError(errors.New("db failed or something")))
			})
		})
	})
})
