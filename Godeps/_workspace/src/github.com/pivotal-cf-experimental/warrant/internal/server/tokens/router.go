package tokens

import (
	"github.com/gorilla/mux"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
)

func NewRouter(tokens *domain.Tokens, users *domain.Users, clients *domain.Clients, publicKey string, urlFinder urlFinder) *mux.Router {
	router := mux.NewRouter()

	router.Handle("/oauth/token", tokenHandler{clients, urlFinder, publicKey}).Methods("POST")
	router.Handle("/oauth/authorize", authorizeHandler{tokens, users}).Methods("POST")
	router.Handle("/token_key", keyHandler{tokens}).Methods("GET")

	return router
}
