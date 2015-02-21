package fakes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/pivotal-cf-experimental/rainmaker/internal/documents"
)

func (fake *CloudController) CreateUser(w http.ResponseWriter, req *http.Request) {
	var document documents.CreateUserRequest
	now := time.Now().UTC()
	err := json.NewDecoder(req.Body).Decode(&document)
	if err != nil {
		panic(err)
	}

	user := NewUser(NewGUID("user"))
	user.GUID = document.GUID
	user.DefaultSpaceGUID = document.DefaultSpaceGUID
	user.CreatedAt = now
	user.UpdatedAt = now

	fake.Users.Add(user)

	response, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
