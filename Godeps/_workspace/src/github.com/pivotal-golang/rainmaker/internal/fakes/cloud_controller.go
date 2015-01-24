package fakes

import (
	"net/http"
	"net/http/httptest"
	"regexp"

	"github.com/gorilla/mux"
)

type CloudController struct {
	server        *httptest.Server
	Organizations *Organizations
	Spaces        *Spaces
	Users         *Users
}

func NewCloudController() *CloudController {
	fake := &CloudController{
		Organizations: NewOrganizations(),
		Spaces:        NewSpaces(),
		Users:         NewUsers(),
	}

	router := mux.NewRouter()
	router.HandleFunc("/v2/organizations", fake.CreateOrganization).Methods("POST")
	router.HandleFunc("/v2/organizations/{guid}", fake.GetOrganization).Methods("GET")
	router.HandleFunc("/v2/organizations/{guid}/users", fake.GetOrganizationUsers).Methods("GET")
	router.HandleFunc("/v2/organizations/{guid}/users/{user_guid}", fake.AssociateUserToOrganization).Methods("PUT")
	router.HandleFunc("/v2/organizations/{guid}/billing_managers", fake.GetOrganizationBillingManagers).Methods("GET")
	router.HandleFunc("/v2/organizations/{guid}/billing_managers/{billing_manager_guid}", fake.AssociateBillingManagerToOrganization).Methods("PUT")
	router.HandleFunc("/v2/organizations/{guid}/auditors", fake.GetOrganizationAuditors).Methods("GET")
	router.HandleFunc("/v2/organizations/{guid}/auditors/{auditor_guid}", fake.AssociateAuditorToOrganization).Methods("PUT")
	router.HandleFunc("/v2/organizations/{guid}/managers", fake.GetOrganizationManagers).Methods("GET")
	router.HandleFunc("/v2/organizations/{guid}/managers/{manager_guid}", fake.AssociateManagerToOrganization).Methods("PUT")
	router.HandleFunc("/v2/spaces", fake.CreateSpace).Methods("POST")
	router.HandleFunc("/v2/spaces/{guid}", fake.GetSpace).Methods("GET")
	router.HandleFunc("/v2/spaces/{guid}/developers/{developer_guid}", fake.AssociateDeveloperToSpace).Methods("PUT")
	router.HandleFunc("/v2/spaces/{guid}/developers", fake.GetSpaceDevelopers).Methods("GET")
	router.HandleFunc("/v2/users", fake.GetUsers).Methods("GET")
	router.HandleFunc("/v2/users", fake.CreateUser).Methods("POST")
	router.HandleFunc("/v2/users/{guid}", fake.GetUser).Methods("GET")

	handler := fake.RequireToken(router)
	fake.server = httptest.NewUnstartedServer(handler)
	return fake
}

func (fake *CloudController) Start() {
	fake.server.Start()
}

func (fake *CloudController) Close() {
	fake.server.Close()
}

func (fake *CloudController) URL() string {
	return fake.server.URL
}

func (fake *CloudController) Reset() {
	fake.Organizations.Clear()
	fake.Spaces.Clear()
	fake.Users.Clear()
}

func (fake *CloudController) RequireToken(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ok, err := regexp.MatchString(`Bearer .+`, req.Header.Get("Authorization"))
		if err != nil {
			panic(err)
		}

		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Not Authorized"))
			return
		}

		handler.ServeHTTP(w, req)
	})
}
