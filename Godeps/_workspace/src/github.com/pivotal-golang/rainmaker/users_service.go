package rainmaker

import (
	"encoding/json"
	"net/http"

	"github.com/pivotal-golang/rainmaker/internal/documents"
)

type UsersService struct {
	config Config
	user   User
}

func NewUsersService(config Config) *UsersService {
	return &UsersService{
		config: config,
	}
}

func (service *UsersService) Create(guid, token string) (User, error) {
	_, body, err := NewClient(service.config).makeRequest(requestArguments{
		Method: "POST",
		Path:   "/v2/users",
		Body: documents.CreateUserRequest{
			GUID: guid,
		},
		Token: token,
		AcceptableStatusCodes: []int{http.StatusCreated},
	})
	if err != nil {
		return User{}, err
	}

	var document documents.UserResponse
	err = json.Unmarshal(body, &document)
	if err != nil {
		panic(err)
	}

	return NewUserFromResponse(document), nil
}

func (service UsersService) Get(guid, token string) (User, error) {
	return FetchUser(service.config, "/v2/users/"+guid, token)
}
