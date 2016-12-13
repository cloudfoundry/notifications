package tokens

import (
	"github.com/gorilla/mux"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
)

func NewRouter(
	tokens *domain.Tokens,
	users *domain.Users,
	clients *domain.Clients,
	publicKey string,
	privateKey string,
	urlFinder urlFinder) *mux.Router {

	router := mux.NewRouter()

	router.Handle("/oauth/token", tokenHandler{clients, users, urlFinder, privateKey}).Methods("POST")
	router.Handle("/oauth/authorize", authorizeHandler{tokens, users, clients}).Methods("POST")
	router.Handle("/token_key", keyHandler{publicKey}).Methods("GET")

	return router
}
