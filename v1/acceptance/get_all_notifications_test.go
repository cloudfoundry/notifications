package v1

import (
	"net/http"

	"github.com/cloudfoundry-incubator/notifications/v1/acceptance/support"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Get a list of all notifications", func() {
	It("allows a user to get body templates", func() {
		client := support.NewClient(Servers.Notifications.URL())

		By("setting the notifications for several clients", func() {
			status, err := client.Notifications.Register(GetClientTokenFor("client-123").Access, support.RegisterClient{
				SourceName: "source name stuff",
				Notifications: map[string]support.RegisterNotification{
					"kind-asd": {
						Description: "remember stuff",
						Critical:    false,
					},
					"kind-abc": {
						Description: "forgot things",
						Critical:    true,
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))

			status, err = client.Notifications.Register(GetClientTokenFor("client-456").Access, support.RegisterClient{
				SourceName: "raptors",
				Notifications: map[string]support.RegisterNotification{
					"dino-kind": {
						Description: "rawr!",
						Critical:    true,
					},
					"fossilized-kind": {
						Description: "crunch!",
						Critical:    false,
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))

			status, err = client.Notifications.Register(GetClientTokenFor("client-890").Access, support.RegisterClient{
				SourceName: "this client has no notifications",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusNoContent))
		})

		By("confirming that all the notifications were registered", func() {
			status, list, err := client.Notifications.List(GetClientTokenFor("notifications-sender").Access)
			Expect(err).NotTo(HaveOccurred())
			Expect(status).To(Equal(http.StatusOK))
			Expect(list).To(HaveLen(3))

			client123 := list["client-123"]
			Expect(client123.Name).To(Equal("source name stuff"))
			Expect(client123.Template).To(Equal("default"))
			Expect(client123.Notifications).To(HaveLen(2))

			kindASD := client123.Notifications["kind-asd"]
			Expect(kindASD.Description).To(Equal("remember stuff"))
			Expect(kindASD.Template).To(Equal("default"))
			Expect(kindASD.Critical).To(BeFalse())

			kindABC := client123.Notifications["kind-abc"]
			Expect(kindABC.Description).To(Equal("forgot things"))
			Expect(kindABC.Template).To(Equal("default"))
			Expect(kindABC.Critical).To(BeTrue())

			client456 := list["client-456"]
			Expect(client456.Name).To(Equal("raptors"))
			Expect(client456.Template).To(Equal("default"))
			Expect(client456.Notifications).To(HaveLen(2))

			dinoKind := client456.Notifications["dino-kind"]
			Expect(dinoKind.Description).To(Equal("rawr!"))
			Expect(dinoKind.Template).To(Equal("default"))
			Expect(dinoKind.Critical).To(BeTrue())

			fossilizedKind := client456.Notifications["fossilized-kind"]
			Expect(fossilizedKind.Description).To(Equal("crunch!"))
			Expect(fossilizedKind.Template).To(Equal("default"))
			Expect(fossilizedKind.Critical).To(BeFalse())

			client890 := list["client-890"]
			Expect(client890.Name).To(Equal("this client has no notifications"))
			Expect(client890.Template).To(Equal("default"))
			Expect(client890.Notifications).To(HaveLen(0))
		})
	})
})
