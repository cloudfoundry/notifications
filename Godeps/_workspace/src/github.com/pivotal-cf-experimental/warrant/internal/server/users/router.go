package users

import (
	"github.com/gorilla/mux"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
)

func NewRouter(users *domain.Users, tokens *domain.Tokens) *mux.Router {
	router := mux.NewRouter()

	router.Handle("/Users", createHandler{users, tokens}).Methods("POST")
	router.Handle("/Users", listHandler{users}).Methods("GET")
	router.Handle("/Users/{guid}", getHandler{users, tokens}).Methods("GET")
	router.Handle("/Users/{guid}", deleteHandler{users, tokens}).Methods("DELETE")
	router.Handle("/Users/{guid}", updateHandler{users, tokens}).Methods("PUT")
	router.Handle("/Users/{guid}/password", passwordHandler{users, tokens}).Methods("PUT")

	return router
}
