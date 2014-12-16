package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/web/handlers"
	"github.com/cloudfoundry-incubator/notifications/web/params"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AssignNotificationTemplate", func() {
	var handler handlers.AssignNotificationTemplate
	var templateAssigner *fakes.TemplateAssigner
	var errorWriter *fakes.ErrorWriter

	BeforeEach(func() {
		templateAssigner = fakes.NewTemplateAssigner()
		errorWriter = fakes.NewErrorWriter()
		handler = handlers.NewAssignNotificationTemplate(templateAssigner, errorWriter)
	})

	It("associates a template with a notification", func() {
		body, err := json.Marshal(map[string]string{
			"template": "my-template",
		})
		if err != nil {
			panic(err)
		}

		w := httptest.NewRecorder()
		request, err := http.NewRequest("PUT", "/clients/my-client/notifications/my-notification/template", bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}

		handler.ServeHTTP(w, request, nil)

		Expect(w.Code).To(Equal(http.StatusNoContent))
		Expect(templateAssigner.AssignToNotificationArguments).To(Equal([]string{"my-client", "my-notification", "my-template"}))
	})

	It("delegates to the error writer when the assigner errors", func() {
		templateAssigner.AssignToNotificationError = errors.New("banana")
		body, err := json.Marshal(map[string]string{
			"template": "my-template",
		})
		if err != nil {
			panic(err)
		}

		w := httptest.NewRecorder()
		request, err := http.NewRequest("PUT", "/clients/my-client/notifications/my-notification/template", bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}

		handler.ServeHTTP(w, request, nil)
		Expect(errorWriter.Error).To(Equal(errors.New("banana")))
	})

	It("writes a ParseError to the error writer when request body is invalid", func() {
		body := []byte(`{ "this is" : not-valid-json }`)

		w := httptest.NewRecorder()
		request, err := http.NewRequest("PUT", "/clients/my-client/notifications/my-notification/template", bytes.NewBuffer(body))
		if err != nil {
			panic(err)
		}

		handler.ServeHTTP(w, request, nil)
		Expect(errorWriter.Error).To(BeAssignableToTypeOf(params.ParseError{}))
	})
})
