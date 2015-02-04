package rainmaker_test

import (
	"net/url"

	"github.com/pivotal-golang/rainmaker"
	"github.com/pivotal-golang/rainmaker/internal/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("UsersList", func() {
	var config rainmaker.Config
	var path, token string
	var list rainmaker.UsersList

	BeforeEach(func() {
		config = rainmaker.Config{
			Host: fakeCloudController.URL(),
		}
		path = "/v2/users"
		token = "token"
		query := url.Values{}
		query.Add("results-per-page", "2")
		query.Add("page", "1")
		list = rainmaker.NewUsersList(config, rainmaker.NewRequestPlan(path, query))
	})

	Describe("Create", func() {
		It("adds a user to the list", func() {
			user, err := list.Create(rainmaker.User{GUID: "user-123"}, token)
			Expect(err).NotTo(HaveOccurred())

			Expect(user.GUID).To(Equal("user-123"))

			err = list.Fetch(token)
			if err != nil {
				panic(err)
			}

			Expect(list.Users).To(HaveLen(1))
			Expect(list.Users[0].GUID).To(Equal("user-123"))
		})
	})

	Describe("Next", func() {
		BeforeEach(func() {
			_, err := list.Create(rainmaker.User{GUID: "user-123"}, token)
			if err != nil {
				panic(err)
			}

			_, err = list.Create(rainmaker.User{GUID: "user-456"}, token)
			if err != nil {
				panic(err)
			}

			_, err = list.Create(rainmaker.User{GUID: "user-789"}, token)
			if err != nil {
				panic(err)
			}
		})

		It("returns the next UserList result for the paginated set", func() {
			err := list.Fetch(token)
			Expect(err).NotTo(HaveOccurred())

			Expect(list.Users).To(HaveLen(2))
			Expect(list.HasNextPage()).To(BeTrue())
			Expect(list.HasPrevPage()).To(BeFalse())
			Expect(list.TotalResults).To(Equal(3))
			Expect(list.TotalPages).To(Equal(2))

			nextList, err := list.Next(token)
			Expect(err).NotTo(HaveOccurred())
			Expect(nextList.Users).To(HaveLen(1))
			Expect(nextList.HasNextPage()).To(BeFalse())
			Expect(nextList.HasPrevPage()).To(BeTrue())
			Expect(nextList.TotalResults).To(Equal(3))
			Expect(nextList.TotalPages).To(Equal(2))

			var users []rainmaker.User
			users = append(users, list.Users...)
			users = append(users, nextList.Users...)
			Expect(users).To(HaveLen(3))

			var guids []string
			for _, user := range users {
				guids = append(guids, user.GUID)
			}
			Expect(guids).To(ConsistOf([]string{"user-123", "user-456", "user-789"}))
		})
	})

	Describe("Prev", func() {
		BeforeEach(func() {
			_, err := list.Create(rainmaker.User{GUID: "user-abc"}, token)
			if err != nil {
				panic(err)
			}

			_, err = list.Create(rainmaker.User{GUID: "user-def"}, token)
			if err != nil {
				panic(err)
			}

			_, err = list.Create(rainmaker.User{GUID: "user-xyz"}, token)
			if err != nil {
				panic(err)
			}
		})

		It("returns the previous UserList result for the paginated set", func() {
			query := url.Values{}
			query.Set("page", "2")
			query.Set("results-per-page", "2")

			list := rainmaker.NewUsersList(config, rainmaker.NewRequestPlan(path, query))
			err := list.Fetch(token)
			if err != nil {
				panic(err)
			}

			Expect(list.Users).To(HaveLen(1))
			Expect(list.HasNextPage()).To(BeFalse())
			Expect(list.HasPrevPage()).To(BeTrue())
			Expect(list.TotalResults).To(Equal(3))
			Expect(list.TotalPages).To(Equal(2))

			prevList, err := list.Prev(token)
			Expect(err).NotTo(HaveOccurred())
			Expect(prevList.Users).To(HaveLen(2))
			Expect(prevList.HasNextPage()).To(BeTrue())
			Expect(prevList.HasPrevPage()).To(BeFalse())
			Expect(prevList.TotalResults).To(Equal(3))
			Expect(prevList.TotalPages).To(Equal(2))

			var users []rainmaker.User
			users = append(users, list.Users...)
			users = append(users, prevList.Users...)
			Expect(users).To(HaveLen(3))

			var guids []string
			for _, user := range users {
				guids = append(guids, user.GUID)
			}
			Expect(guids).To(ConsistOf([]string{"user-abc", "user-def", "user-xyz"}))
		})
	})

	Describe("HasNextPage", func() {
		It("indicates whether or not there is a next page of results", func() {
			list.NextURL = "/v2/users?page=2"
			Expect(list.HasNextPage()).To(BeTrue())

			list.NextURL = ""
			Expect(list.HasNextPage()).To(BeFalse())
		})
	})

	Describe("HasPrevPage", func() {
		It("indicates whether or not there is a previous page of results", func() {
			list.PrevURL = "/v2/users?page=1"
			Expect(list.HasPrevPage()).To(BeTrue())

			list.PrevURL = ""
			Expect(list.HasPrevPage()).To(BeFalse())
		})
	})

	Describe("AllUsers", func() {
		BeforeEach(func() {
			_, err := list.Create(rainmaker.User{GUID: "user-abc"}, token)
			if err != nil {
				panic(err)
			}

			_, err = list.Create(rainmaker.User{GUID: "user-def"}, token)
			if err != nil {
				panic(err)
			}

			_, err = list.Create(rainmaker.User{GUID: "user-xyz"}, token)
			if err != nil {
				panic(err)
			}
		})

		It("returns a slice of all of users", func() {
			err := list.Fetch(token)
			Expect(err).NotTo(HaveOccurred())

			users, err := list.AllUsers(token)
			Expect(err).NotTo(HaveOccurred())

			Expect(users).To(HaveLen(3))
			var guids []string
			for _, user := range users {
				guids = append(guids, user.GUID)
			}
			Expect(guids).To(ConsistOf([]string{"user-abc", "user-def", "user-xyz"}))
		})
	})

	Describe("Associate", func() {
		It("associates a user with the listed resource", func() {
			spaceGUID := "space-abc"
			fakeCloudController.Spaces.Add(fakes.Space{
				GUID:       spaceGUID,
				Developers: fakes.NewUsers(),
			})

			user, err := list.Create(rainmaker.User{GUID: "user-abc"}, token)
			if err != nil {
				panic(err)
			}

			list = rainmaker.NewUsersList(config, rainmaker.NewRequestPlan("/v2/spaces/"+spaceGUID+"/developers", url.Values{}))
			err = list.Associate(user.GUID, token)
			Expect(err).NotTo(HaveOccurred())

			err = list.Fetch(token)
			Expect(err).NotTo(HaveOccurred())
			Expect(list.Users).To(HaveLen(1))
			Expect(list.Users[0].GUID).To(Equal(user.GUID))
		})
	})
})
