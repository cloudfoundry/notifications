package tokens

import (
	"encoding/json"
	"net/http"

	"github.com/pivotal-cf-experimental/warrant/internal/documents"
	"github.com/pivotal-cf-experimental/warrant/internal/server/domain"
)

type keyHandler struct {
	tokens *domain.Tokens
}

func (h keyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	response, err := json.Marshal(documents.TokenKeyResponse{
		Alg:   "SHA256withRSA",
		Value: h.tokens.PublicKey,
		Kty:   "RSA",
		Use:   "sig",
		N:     "ANJufZdrvYg5zG61x36pDq59nVUN73wSanA7hVCtN3ftT2Rm1ZTQqp5KSCfLMhaaVvJY51sHj+/i4lqUaM9CO32G93fE44VfOmPfexZeAwa8YDOikyTrhP7sZ6A4WUNeC4DlNnJF4zsznU7JxjCkASwpdL6XFwbRSzGkm6b9aM4vIewyclWehJxUGVFhnYEzIQ65qnr38feVP9enOVgQzpKsCJ+xpa8vZ/UrscoG3/IOQM6VnLrGYAyyCGeyU1JXQW/KlNmtA5eJry2Tp+MD6I34/QsNkCArHOfj8H9tXz/oc3/tVkkR252L/Lmp0TtIGfHpBmoITP9h+oKiW6NpyCc=",
		E:     "AQAB",
	})
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
