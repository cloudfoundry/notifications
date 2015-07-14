package clients

import (
	"github.com/gorilla/mux"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
)

func NewRouter(clients *domain.Clients, tokens *domain.Tokens) *mux.Router {
	router := mux.NewRouter()

	router.Handle("/oauth/clients", createHandler{clients, tokens}).Methods("POST")
	router.Handle("/oauth/clients/{guid}", getHandler{clients, tokens}).Methods("GET")
	router.Handle("/oauth/clients/{guid}", deleteHandler{clients, tokens}).Methods("DELETE")

	return router
}
