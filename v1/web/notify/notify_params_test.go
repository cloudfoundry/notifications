package notify_test

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/cloudfoundry-incubator/notifications/v1/web/notify"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NotifyParams", func() {
	Describe("NewNotifyParams", func() {
		It("parses the body of the given request", func() {
			parameters, err := notify.NewNotifyParams(ioutil.NopCloser(strings.NewReader(`{
                "kind_id": "test_email",
                "reply_to": "me@awesome.com",
                "subject": "Summary of contents",
                "text": "Contents of the email message"
            }`)))
			Expect(err).NotTo(HaveOccurred())

			Expect(parameters.KindID).To(Equal("test_email"))
			Expect(parameters.KindDescription).To(Equal(""))
			Expect(parameters.SourceDescription).To(Equal(""))
			Expect(parameters.ReplyTo).To(Equal("me@awesome.com"))
			Expect(parameters.Subject).To(Equal("Summary of contents"))
			Expect(parameters.Text).To(Equal("Contents of the email message"))
		})

		It("does not blow up if the request body is empty", func() {
			Expect(func() {
				notify.NewNotifyParams(ioutil.NopCloser(strings.NewReader("")))
			}).NotTo(Panic())
		})

		Describe("to field parsing", func() {
			It("handles when a name is attached to the address", func() {
				parameters, err := notify.NewNotifyParams(ioutil.NopCloser(strings.NewReader(`{
					"to": "The User <user@example.com>"
				}`)))
				Expect(err).NotTo(HaveOccurred())
				Expect(parameters.To).To(Equal("user@example.com"))
			})

			It("populates the To field with the parsed email address", func() {
				parameters, err := notify.NewNotifyParams(ioutil.NopCloser(strings.NewReader(`{
                    "to": "user@example.com"
				}`)))
				Expect(err).NotTo(HaveOccurred())
				Expect(parameters.To).To(Equal("user@example.com"))
			})

			It("sets the to field to InvalidEmail cannot be parsed", func() {
				parameters, err := notify.NewNotifyParams(ioutil.NopCloser(strings.NewReader(`{
                    "to": "<The User"
				}`)))
				Expect(err).NotTo(HaveOccurred())
				Expect(parameters.To).To(Equal(notify.InvalidEmail))
			})

			It("Sets the To field to empty of if it is not specified", func() {
				parameters, err := notify.NewNotifyParams(ioutil.NopCloser(strings.NewReader(`{
                    "to": ""
				}`)))
				Expect(err).NotTo(HaveOccurred())
				Expect(parameters.To).To(Equal(""))
			})
		})

		Describe("role field parsing", func() {
			It("sets the role field to empty if it is not specificed", func() {
				parameters, err := notify.NewNotifyParams(ioutil.NopCloser(strings.NewReader("{}")))
				Expect(err).NotTo(HaveOccurred())
				Expect(parameters.Role).To(Equal(""))

			})

			It("sets the role field that is specified", func() {
				parameters, err := notify.NewNotifyParams(ioutil.NopCloser(strings.NewReader(`{
                    "role": "the-role"
				}`)))
				Expect(err).NotTo(HaveOccurred())
				Expect(parameters.Role).To(Equal("the-role"))
			})
		})

		Describe("html parsing", func() {
			Context("when a doctype is passed in", func() {
				It("pulls out the doctype", func() {
					parameters, err := notify.NewNotifyParams(ioutil.NopCloser(strings.NewReader(`{
                        "kind_id": "test_email",
                        "html": "<!DOCTYPE html>"
					}`)))
					Expect(err).NotTo(HaveOccurred())
					Expect(parameters.ParsedHTML.Doctype).To(Equal("<!DOCTYPE html>"))
				})
			})

			Context("when no doctype is passed", func() {
				It("returns an empty doctype", func() {
					parameters, err := notify.NewNotifyParams(ioutil.NopCloser(strings.NewReader(`{
                        "kind_id": "test_email",
                        "html": ""
					}`)))
					Expect(err).NotTo(HaveOccurred())
					Expect(parameters.ParsedHTML.Doctype).To(Equal(""))
				})
			})

			Context("when a head tag is passed in", func() {
				It("pulls out the contents of the head tag", func() {
					parameters, err := notify.NewNotifyParams(ioutil.NopCloser(strings.NewReader(`{
                        "kind_id": "test_email",
                        "html": "<head><title>BananaDamage</title></head>"
					}`)))
					Expect(err).NotTo(HaveOccurred())
					Expect(parameters.ParsedHTML.Head).To(Equal("<title>BananaDamage</title>"))
				})
			})

			Context("when no head tag is passed in", func() {
				It("Head is left as an empty string", func() {
					parameters, err := notify.NewNotifyParams(ioutil.NopCloser(strings.NewReader(`{
                        "kind_id": "test_email",
                        "html": ""
					}`)))
					Expect(err).NotTo(HaveOccurred())
					Expect(parameters.ParsedHTML.Head).To(Equal(""))
				})
			})

			Context("body tags are present", func() {
				var body io.ReadCloser

				BeforeEach(func() {
					body = ioutil.NopCloser(strings.NewReader(`{
                        "kind_id": "test_email",
                        "html": "<body class='bananaDamage'><p>The TEXT</p><h1>the TITLE</h1></body>"
                    }`))
				})

				It("pulls out the html in the body", func() {
					parameters, err := notify.NewNotifyParams(body)
					Expect(err).NotTo(HaveOccurred())
					Expect(parameters.ParsedHTML.BodyContent).To(ContainSubstring("<p>The TEXT</p><h1>the TITLE</h1>"))
				})

				It("preserves any attributes on the body tag itself", func() {
					parameters, err := notify.NewNotifyParams(body)
					Expect(err).NotTo(HaveOccurred())
					Expect(parameters.ParsedHTML.BodyAttributes).To(ContainSubstring(`class="bananaDamage"`))
				})
			})

			Context("when only an html tag is present", func() {
				It("the contents in the html tag are put into the body", func() {
					parameters, err := notify.NewNotifyParams(ioutil.NopCloser(strings.NewReader(`{
                        "kind_id": "test_email",
                        "html": "<html><head><title>BananaDamage</title></head><p>The TEXT</p><h1>the TITLE</h1></html>"
                    }`)))
					Expect(err).NotTo(HaveOccurred())
					Expect(parameters.ParsedHTML.BodyContent).To(Equal("<p>The TEXT</p><h1>the TITLE</h1>"))
					Expect(parameters.ParsedHTML.Head).To(Equal("<title>BananaDamage</title>"))
				})
			})

			Context("when just bare html is passed without surrounding html/body tags", func() {
				It("the html is placed in the body", func() {
					parameters, err := notify.NewNotifyParams(ioutil.NopCloser(strings.NewReader(`{
                        "kind_id": "test_email",
                        "html": "<p>The TEXT</p><h1>the TITLE</h1>"
                    }`)))
					Expect(err).NotTo(HaveOccurred())
					Expect(parameters.ParsedHTML.BodyContent).To(Equal("<p>The TEXT</p><h1>the TITLE</h1>"))
				})
			})

			Context("when invalid html is passed", func() {
				It("pulls out the html anyway", func() {
					parameters, err := notify.NewNotifyParams(ioutil.NopCloser(strings.NewReader(`{
                        "kind_id": "test_email",
                        "html": "<html><p>The TEXT<h1>the TITLE</h1></html>"
                    }`)))
					Expect(err).NotTo(HaveOccurred())
					Expect(parameters.ParsedHTML.BodyContent).To(Equal("<p>The TEXT</p><h1>the TITLE</h1>"))

					parameters, err = notify.NewNotifyParams(ioutil.NopCloser(strings.NewReader(`{
                        "kind_id": "test_email",
                        "html": "<html><p>The TEXT<h1>the TITLE</h1></body>"
                    }`)))
					Expect(err).NotTo(HaveOccurred())
					Expect(parameters.ParsedHTML.BodyContent).To(Equal("<p>The TEXT</p><h1>the TITLE</h1>"))
				})
			})

			Context("when no html is passed", func() {
				It("does not error", func() {
					parameters, err := notify.NewNotifyParams(ioutil.NopCloser(strings.NewReader(`{
                        "kind_id": "test_email",
                        "text": "not html yo"
                    }`)))
					Expect(err).NotTo(HaveOccurred())
					Expect(parameters.ParsedHTML.BodyContent).To(Equal(""))
				})
			})

			Context("when the to field is invalid", func() {
				Context("when it has unmatched <", func() {
					It("assigns <>invalidEmail<>", func() {
						parameters, err := notify.NewNotifyParams(ioutil.NopCloser(strings.NewReader(`{
							"to": "<invalid email",
							"text": "Contents of the email message"
						}`)))
						Expect(err).NotTo(HaveOccurred())
						Expect(parameters.To).To(Equal(notify.InvalidEmail))
					})
				})

				Context("when it is missing an @", func() {
					It("assigns <>invalidEmail<>", func() {
						parameters, err := notify.NewNotifyParams(ioutil.NopCloser(strings.NewReader(`{
							"to": "invalidemail.com",
							"text": "Contents of the email message"
						}`)))
						Expect(err).NotTo(HaveOccurred())
						Expect(parameters.To).To(Equal(notify.InvalidEmail))
					})
				})
			})

			Context("when a lot of complicated html is sent", func() {
				It("does the right thing", func() {
					html := `<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.0 Transitional//EN\"><head><title>New Relic</title></head><body bgcolor=\"#cccccc\" leftmargin=\"10\" topmargin=\"0\" rightmargin=\"10\" bottommargin=\"10\" marginheight=\"10\" marginwidth=\"10\"><div>div here ya</div></body>`
					body := ioutil.NopCloser(strings.NewReader(`{"kind_id": "test_email", "html": "` + html + `"}`))

					parameters, err := notify.NewNotifyParams(body)
					Expect(err).NotTo(HaveOccurred())
					Expect(parameters.ParsedHTML.Doctype).To(Equal("<!DOCTYPE HTML PUBLIC \"-//W3C//DTD HTML 4.0 Transitional//EN\">"))
					Expect(parameters.ParsedHTML.BodyAttributes).To(Equal("bgcolor=\"#cccccc\" leftmargin=\"10\" topmargin=\"0\" rightmargin=\"10\" bottommargin=\"10\" marginheight=\"10\" marginwidth=\"10\""))
					Expect(parameters.ParsedHTML.BodyContent).To(Equal("<div>div here ya</div>"))
					Expect(parameters.ParsedHTML.Head).To(Equal("<title>New Relic</title>"))
				})
			})
		})
	})
})
