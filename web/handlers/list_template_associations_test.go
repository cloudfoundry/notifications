package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TemplateAssociation struct {
	Client       string
	Notification string
}

var _ = Describe("ListTemplateAssociations", func() {
	var handler handlers.ListTemplateAssociations
	var templateID string
	var writer *httptest.ResponseRecorder
	var request *http.Request
	var lister *fakes.TemplateAssociationLister
	var errorWriter *fakes.ErrorWriter

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
		handler = handlers.NewListTemplateAssociations(lister, errorWriter)

		writer = httptest.NewRecorder()
		request, err = http.NewRequest("GET", "/templates/"+templateID+"/associations", nil)
		if err != nil {
			panic(err)
		}
	})

	It("returns a list of clients and notifications associated to the given template", func() {
		handler.ServeHTTP(writer, request, nil)

		Expect(writer.Code).To(Equal(http.StatusOK))

		var assoc struct {
			Associations []TemplateAssociation
		}
		err := json.Unmarshal(writer.Body.Bytes(), &assoc)
		if err != nil {
			panic(err)
		}
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
	})

	Context("when errors occur", func() {
		Context("when the lister service returns an error", func() {
			It("delegates to the error handler", func() {
				lister.ListError = errors.New("db failed or something")

				handler.ServeHTTP(writer, request, nil)
				Expect(errorWriter.Error).To(MatchError(errors.New("db failed or something")))
			})
		})
	})
})
