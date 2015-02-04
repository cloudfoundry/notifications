package rainmaker

import (
	"encoding/json"
	"net/http"
	"net/url"
	"path"

	"github.com/pivotal-golang/rainmaker/internal/documents"
)

type UsersList struct {
	config       Config
	plan         requestPlan
	TotalResults int
	TotalPages   int
	NextURL      string
	PrevURL      string
	Users        []User
}

func NewUsersList(config Config, plan requestPlan) UsersList {
	return UsersList{
		config: config,
		plan:   plan,
	}
}

func newUsersListFromResponse(config Config, plan requestPlan, response documents.UsersListResponse) UsersList {
	list := NewUsersList(config, plan)
	list.TotalResults = response.TotalResults
	list.TotalPages = response.TotalPages
	list.PrevURL = response.PrevURL
	list.NextURL = response.NextURL
	list.Users = make([]User, 0)

	for _, userResponse := range response.Resources {
		list.Users = append(list.Users, newUserFromResponse(config, userResponse))
	}

	return list
}

func (list UsersList) HasNextPage() bool {
	return list.NextURL != ""
}

func (list UsersList) HasPrevPage() bool {
	return list.PrevURL != ""
}

func (list UsersList) Next(token string) (UsersList, error) {
	nextURL, err := url.Parse("http://example.com" + list.NextURL)
	if err != nil {
		return UsersList{}, err
	}

	nextList := NewUsersList(list.config, newRequestPlan(nextURL.Path, nextURL.Query()))
	err = nextList.Fetch(token)

	return nextList, err
}

func (list UsersList) Prev(token string) (UsersList, error) {
	prevURL, err := url.Parse("http://example.com" + list.PrevURL)
	if err != nil {
		return UsersList{}, err
	}

	prevList := NewUsersList(list.config, newRequestPlan(prevURL.Path, prevURL.Query()))
	err = prevList.Fetch(token)

	return prevList, err
}

func (list UsersList) AllUsers(token string) ([]User, error) {
	l := list
	users := make([]User, 0)

	for l.HasPrevPage() {
		var err error
		l, err = l.Prev(token)
		if err != nil {
			return []User{}, err
		}

		users = append(l.Users, users...)
	}

	users = append(users, list.Users...)

	for l.HasNextPage() {
		var err error
		l, err = l.Next(token)
		if err != nil {
			return []User{}, err
		}

		users = append(users, l.Users...)
	}

	return users, nil
}

func (list UsersList) Create(user User, token string) (User, error) {
	var document documents.UserResponse

	_, body, err := NewClient(list.config).makeRequest(requestArguments{
		Method: "POST",
		Path:   list.plan.Path,
		Token:  token,
		Body:   user,
		AcceptableStatusCodes: []int{http.StatusCreated},
	})
	if err != nil {
		return User{}, err
	}

	err = json.Unmarshal(body, &document)
	if err != nil {
		panic(err)
	}

	return newUserFromResponse(list.config, document), nil
}

func (list UsersList) Associate(userGUID, token string) error {
	_, _, err := NewClient(list.config).makeRequest(requestArguments{
		Method: "PUT",
		Path:   path.Join(list.plan.Path, userGUID),
		Token:  token,
		AcceptableStatusCodes: []int{http.StatusCreated},
	})

	return err
}

func (list *UsersList) Fetch(token string) error {
	_, body, err := NewClient(list.config).makeRequest(requestArguments{
		Method: "GET",
		Path:   list.plan.Path,
		Query:  list.plan.Query,
		Token:  token,
		AcceptableStatusCodes: []int{http.StatusOK},
	})
	if err != nil {
		return err
	}

	var response documents.UsersListResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	updatedList := newUsersListFromResponse(list.config, list.plan, response)
	list.TotalResults = updatedList.TotalResults
	list.TotalPages = updatedList.TotalPages
	list.NextURL = updatedList.NextURL
	list.PrevURL = updatedList.PrevURL
	list.Users = updatedList.Users

	return nil
}
