package warrant_test

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/pivotal-cf-experimental/warrant"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GroupsService", func() {
	var (
		service warrant.GroupsService
		token   string
		config  warrant.Config
	)

	BeforeEach(func() {
		config = warrant.Config{
			Host:          fakeUAA.URL(),
			SkipVerifySSL: true,
			TraceWriter:   TraceWriter,
		}
		service = warrant.NewGroupsService(config)
		token = fakeUAA.ClientTokenFor("admin", []string{"scim.write", "scim.read"}, []string{"scim"})
	})

	Describe("Create", func() {
		It("creates a group given a name", func() {
			group, err := service.Create("banana.write", token)
			Expect(err).NotTo(HaveOccurred())
			Expect(group.ID).NotTo(BeEmpty())
			Expect(group.DisplayName).To(Equal("banana.write"))
			Expect(group.Version).To(Equal(0))
			Expect(group.CreatedAt).To(BeTemporally("~", time.Now().UTC(), 2*time.Millisecond))
			Expect(group.UpdatedAt).To(BeTemporally("~", time.Now().UTC(), 2*time.Millisecond))
		})

		It("requires the scim.write scope", func() {
			token = fakeUAA.ClientTokenFor("admin", []string{"scim.read"}, []string{"scim"})
			_, err := service.Create("banana.write", token)
			Expect(err).To(BeAssignableToTypeOf(warrant.UnauthorizedError{}))
		})

		It("requires the scim audience", func() {
			token = fakeUAA.ClientTokenFor("admin", []string{"scim.write"}, []string{"banana"})
			_, err := service.Create("banana.write", token)
			Expect(err).To(BeAssignableToTypeOf(warrant.UnauthorizedError{}))
		})

		Context("failure cases", func() {
			It("returns an error when a group with the given name already exists", func() {
				_, err := service.Create("banana.write", token)
				Expect(err).NotTo(HaveOccurred())

				_, err = service.Create("banana.write", token)
				Expect(err).To(BeAssignableToTypeOf(warrant.DuplicateResourceError{}))
				Expect(err.Error()).To(Equal("duplicate resource: {\"message\":\"A group with displayName: banana.write already exists.\",\"error\":\"scim_resource_already_exists\"}"))
			})

			It("returns an error when the json response is malformed", func() {
				malformedJSONServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte("this is not JSON"))
				}))
				service = warrant.NewGroupsService(warrant.Config{
					Host:          malformedJSONServer.URL,
					SkipVerifySSL: true,
					TraceWriter:   TraceWriter,
				})

				_, err := service.Create("banana.read", "some-token")
				Expect(err).To(BeAssignableToTypeOf(warrant.MalformedResponseError{}))
				Expect(err).To(MatchError("malformed response: invalid character 'h' in literal true (expecting 'r')"))
			})
		})
	})

	Describe("Get", func() {
		var createdGroup warrant.Group

		BeforeEach(func() {
			var err error
			createdGroup, err = service.Create("banana.write", token)
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns the found group", func() {
			group, err := service.Get(createdGroup.ID, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(group).To(Equal(createdGroup))
		})

		It("requires the scim.read scope", func() {
			token = fakeUAA.ClientTokenFor("admin", []string{"scim.write"}, []string{"scim"})
			_, err := service.Get(createdGroup.ID, token)
			Expect(err).To(BeAssignableToTypeOf(warrant.UnauthorizedError{}))
		})

		It("requires the scim audience", func() {
			token = fakeUAA.ClientTokenFor("admin", []string{"scim.read"}, []string{"banana"})
			_, err := service.Get(createdGroup.ID, token)
			Expect(err).To(BeAssignableToTypeOf(warrant.UnauthorizedError{}))
		})

		Context("failure cases", func() {
			It("returns an error when the group cannot be found", func() {
				_, err := service.Get("non-existent-group-id", token)
				Expect(err).To(BeAssignableToTypeOf(warrant.NotFoundError{}))
			})

			It("returns an error when the json response is malformed", func() {
				malformedJSONServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Write([]byte("this is not JSON"))
				}))
				service = warrant.NewGroupsService(warrant.Config{
					Host:          malformedJSONServer.URL,
					SkipVerifySSL: true,
					TraceWriter:   TraceWriter,
				})

				_, err := service.Get("some-group-id", "some-token")
				Expect(err).To(BeAssignableToTypeOf(warrant.MalformedResponseError{}))

				Expect(err).To(MatchError("malformed response: invalid character 'h' in literal true (expecting 'r')"))
			})
		})
	})

	Describe("Delete", func() {
		var group warrant.Group

		BeforeEach(func() {
			var err error
			group, err = service.Create("banana.read", token)
			Expect(err).NotTo(HaveOccurred())
		})

		It("deletes the group", func() {
			err := service.Delete(group.ID, token)
			Expect(err).NotTo(HaveOccurred())

			_, err = service.Create("banana.read", token)
			Expect(err).NotTo(HaveOccurred())
		})

		It("requires the scim.write scope", func() {
			token = fakeUAA.ClientTokenFor("admin", []string{"scim.read"}, []string{"scim"})
			err := service.Delete(group.ID, token)
			Expect(err).To(BeAssignableToTypeOf(warrant.UnauthorizedError{}))
		})

		It("requires the scim audience", func() {
			token = fakeUAA.ClientTokenFor("admin", []string{"scim.write"}, []string{"banana"})
			err := service.Delete(group.ID, token)
			Expect(err).To(BeAssignableToTypeOf(warrant.UnauthorizedError{}))
		})

		It("returns an error when the group does not exist", func() {
			err := service.Delete("non-existant-group-guid", token)
			Expect(err).To(BeAssignableToTypeOf(warrant.NotFoundError{}))
		})
	})

	Describe("List", func() {
		It("retrieves a list of all the groups", func() {
			writeGroup, err := service.Create("banana.write", token)
			Expect(err).NotTo(HaveOccurred())

			readGroup, err := service.Create("banana.read", token)
			Expect(err).NotTo(HaveOccurred())

			groups, err := service.List(warrant.Query{}, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(groups).To(HaveLen(2))
			Expect(groups).To(ConsistOf(writeGroup, readGroup))
		})

		Context("failure cases", func() {
			It("returns an error when the server does not respond validly", func() {
				erroringServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
				service = warrant.NewGroupsService(warrant.Config{
					Host:          erroringServer.URL,
					SkipVerifySSL: true,
					TraceWriter:   TraceWriter,
				})

				_, err := service.List(warrant.Query{}, token)
				Expect(err).To(BeAssignableToTypeOf(warrant.UnexpectedStatusError{}))
			})

			It("returns an error when the JSON is malformed", func() {
				malformedJSONServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
					w.Write([]byte("this is not JSON"))
				}))
				service = warrant.NewGroupsService(warrant.Config{
					Host:          malformedJSONServer.URL,
					SkipVerifySSL: true,
					TraceWriter:   TraceWriter,
				})

				_, err := service.List(warrant.Query{}, token)
				Expect(err).To(BeAssignableToTypeOf(warrant.MalformedResponseError{}))
				Expect(err).To(MatchError("malformed response: invalid character 'h' in literal true (expecting 'r')"))
			})
		})
	})
})
