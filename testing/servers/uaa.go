package servers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

type UAA struct {
	server    *httptest.Server
	ServerURL string
}

func NewUAA() *UAA {
	router := mux.NewRouter()
	router.HandleFunc("/oauth/token", UAAPostOAuthToken).Methods("POST")
	router.HandleFunc("/token_key", UAAGetTokenKey).Methods("GET")
	router.HandleFunc("/token_keys", UAAGetTokenKeys).Methods("GET")
	router.Path("/Users").Queries("filter", "").Handler(UAAGetUserFilter).Methods("GET")
	router.HandleFunc("/Users", UAAGetUsers).Methods("GET")
	router.HandleFunc("/Users/{userGUID}", UAAGetUser).Methods("GET")
	router.HandleFunc("/Groups", UAAGetUsersByScope).Methods("GET")
	router.HandleFunc("/{anything:.*}", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("UAA ROUTE REQUEST ---> %+v\n", req)
		w.WriteHeader(http.StatusTeapot)
	}))
	return &UAA{
		server: httptest.NewUnstartedServer(router),
	}
}

func (s *UAA) Boot() {
	s.server.Start()
	s.ServerURL = s.server.URL
	os.Setenv("UAA_HOST", s.server.URL)
}

func (s UAA) Close() {
	s.server.Close()
}

func ReadFile(filename string) string {
	env, err := application.NewEnvironment()
	if err != nil {
		panic(err)
	}

	root := env.RootPath
	fileContents, err := ioutil.ReadFile(root + filename)
	if err != nil {
		panic(err)
	}

	return string(fileContents)
}

var UAAPostOAuthToken = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		panic(err)
	}

	encodedCredentials := strings.TrimPrefix(req.Header.Get("Authorization"), "Basic ")
	decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
	credentialsParts := strings.Split(string(decodedCredentials), ":")
	clientID := credentialsParts[0]

	token := jwt.New(jwt.GetSigningMethod("RS256"))
	claims := jwt.MapClaims{}
	token.Claims = claims
	token.Header["kid"] = "legacy-key-id"
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	claims["client_id"] = clientID
	claims["iss"] = "http://" + req.Host + "/oauth/token"

	switch req.Form.Get("grant_type") {
	case "client_credentials":
		switch clientID {
		case "non-critical-client":
			claims["scope"] = []string{"notifications.manage", "notifications.write",
				"emails.write", "notification_preferences.admin", "notification_templates.admin",
				"notification_templates.write", "notification_templates.read"}
		case "unauthorized-client":
			claims["scope"] = []string{}
		case "non-admin-client":
			claims["scope"] = []string{"notifications.write"}
		case "admin-client":
			claims["scope"] = []string{"notifications.admin"}
		default:
			claims["scope"] = []string{"notifications.manage",
				"notifications.write", "emails.write", "notification_preferences.admin",
				"critical_notifications.write", "notification_templates.admin",
				"notification_templates.write", "notification_templates.read"}
		}
	case "authorization_code":
		claims["user_id"] = strings.TrimSuffix(req.Form.Get("code"), "-code")
		switch claims["user_id"] {
		case "unauthorized-user":
			claims["scope"] = []string{}
		default:
			claims["scope"] = []string{"notification_preferences.read", "notification_preferences.write"}
		}
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(ReadFile("/testing/fixtures/private.pem")))
	if err != nil {
		panic(err)
	}

	tokenString, err := token.SignedString(key)
	if err != nil {
		panic(err)
	}

	response, err := json.Marshal(map[string]string{
		"access_token": tokenString,
	})
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
})

var UAAGetTokenKey = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	response, err := json.Marshal(map[string]string{
		"alg":   "SHA256withRSA",
		"kid":   "legacy-key-id",
		"value": ReadFile("/testing/fixtures/public.pem"),
	})
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
})

var UAAGetTokenKeys = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	response, err := json.Marshal(map[string]interface{}{
		"keys": []interface{}{
			map[string]string{
				"alg":   "SHA256withRSA",
				"kid":   "legacy-key-id",
				"value": ReadFile("/testing/fixtures/public.pem"),
			},
			map[string]string{
				"alg":   "SHA256withRSA",
				"kid":   "new-key-id",
				"value": ReadFile("/testing/fixtures/public.pem"),
			},
		},
	})
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
})

var UAAGetUserFilter = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		panic(err)
	}

	filter := req.Form.Get("filter")
	filterParts := strings.Split(filter, " or ")
	queryRegexp := regexp.MustCompile(`Id eq "(.*)"`)
	resources := []interface{}{}
	for _, part := range filterParts {
		matches := queryRegexp.FindAllStringSubmatch(part, 1)
		match := matches[0]

		if match[1] == "user-malformed-email" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{}`))
			return
		}

		if user, ok := UAAUsers[match[1]]; ok {
			resources = append(resources, user)
		}
	}

	response, err := json.Marshal(map[string]interface{}{
		"resources":    resources,
		"startIndex":   1,
		"itemsPerPage": 100,
		"totalResults": 1,
		"schemas":      []string{"urn:scim:schemas:core:1.0"},
	})
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
})

var UAAGetUser = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	userGUID := mux.Vars(req)["userGUID"]

	user, ok := UAAUsers[userGUID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response, err := json.Marshal(user)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
})

var UAAGetUsers = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	response, err := json.Marshal(map[string]interface{}{
		"resources":    allUsersResponse,
		"startIndex":   1,
		"itemsPerPage": 100,
		"totalResults": 1,
		"schemas":      []string{"urn:scim:schemas:core:1.0"},
	})
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
})

var UAAGetUsersByScope = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		panic(err)
	}

	attribute := req.Form.Get("attributes")
	if attribute != "members" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors": "attribute not found"}`))
	}

	filter := req.Form.Get("filter")
	queryRegexp := regexp.MustCompile(`displayName eq "(.*)"`)
	matches := queryRegexp.FindAllStringSubmatch(filter, 1)
	match := matches[0]
	resources := UAAUsersByScope[match[1]]

	response, err := json.Marshal(map[string]interface{}{
		"resources":    resources,
		"startIndex":   1,
		"itemsPerPage": 100,
		"totalResults": 1,
		"schemas":      []string{"urn:scim:schemas:core:1.0"},
	})
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
})

var UAAUsersByScope = map[string]interface{}{
	"this.scope": []map[string]interface{}{
		{
			"members": []map[string]string{
				{
					"origin": "uaa",
					"type":   "user",
					"value":  "user-369",
				},
			},
		},
	},
}

var UAAUsers = map[string]map[string]interface{}{
	"user-111": {
		"id": "user-111",
		"meta": map[string]interface{}{
			"version":      4,
			"created":      "2014-07-16T21:00:09.021Z",
			"lastModified": "2014-08-04T19:16:29.172Z",
		},
		"userName": "User111",
		"name":     map[string]string{},
		"emails": []map[string]string{
			{"value": "user-111@example.com"},
		},
		"groups": []map[string]string{
			{
				"value":   "some-group-guid",
				"display": "notifications.write",
				"type":    "DIRECT",
			},
		},
		"approvals": []interface{}{},
		"active":    true,
		"verified":  false,
		"origin":    "uaa",
		"schemas":   []string{"urn:scim:schemas:core:1.0"},
	},
	"user-123": {
		"id": "user-123",
		"meta": map[string]interface{}{
			"version":      4,
			"created":      "2014-07-16T21:00:09.021Z",
			"lastModified": "2014-08-04T19:16:29.172Z",
		},
		"userName": "User123",
		"name":     map[string]string{},
		"emails": []map[string]string{
			{"value": "user-123@example.com"},
		},
		"groups": []map[string]string{
			{
				"value":   "some-group-guid",
				"display": "notifications.write",
				"type":    "DIRECT",
			},
		},
		"approvals": []interface{}{},
		"active":    true,
		"verified":  false,
		"origin":    "uaa",
		"schemas":   []string{"urn:scim:schemas:core:1.0"},
	},
	"user-456": {
		"id": "user-456",
		"meta": map[string]interface{}{
			"version":      4,
			"created":      "2014-07-16T21:00:09.021Z",
			"lastModified": "2014-08-04T19:16:29.172Z",
		},
		"userName": "User456",
		"name":     map[string]string{},
		"emails": []map[string]string{
			{"value": "user-456@example.com"},
		},
		"groups": []map[string]string{
			{
				"value":   "some-group-guid",
				"display": "notifications.write",
				"type":    "DIRECT",
			},
		},
		"approvals": []interface{}{},
		"active":    true,
		"verified":  false,
		"origin":    "uaa",
		"schemas":   []string{"urn:scim:schemas:core:1.0"},
	},
	"user-789": {
		"id": "user-789",
		"meta": map[string]interface{}{
			"version":      4,
			"created":      "2014-07-16T21:00:09.021Z",
			"lastModified": "2014-08-04T19:16:29.172Z",
		},
		"userName": "User789",
		"name":     map[string]string{},
		"emails": []map[string]string{
			{"value": "user-789"},
		},
		"groups": []map[string]string{
			{
				"value":   "some-group-guid",
				"display": "notifications.write",
				"type":    "DIRECT",
			},
		},
		"approvals": []interface{}{},
		"active":    true,
		"verified":  false,
		"origin":    "uaa",
		"schemas":   []string{"urn:scim:schemas:core:1.0"},
	},
	"user-369": {
		"id": "user-369",
		"meta": map[string]interface{}{
			"version":      4,
			"created":      "2014-07-16T21:00:09.021Z",
			"lastModified": "2014-08-04T19:16:29.172Z",
		},
		"userName": "User369",
		"name":     map[string]string{},
		"emails": []map[string]string{
			{"value": "user-369@example.com"},
		},
		"groups": []map[string]string{
			{
				"value":   "some-group-guid",
				"display": "notifications.write",
				"type":    "DIRECT",
			},
			{
				"value":   "this-scope-guid",
				"display": "this.scope",
				"type":    "DIRECT",
			},
		},
		"approvals": []interface{}{},
		"active":    true,
		"verified":  false,
		"origin":    "uaa",
		"schemas":   []string{"urn:scim:schemas:core:1.0"},
	},
	"091b6583-0933-4d17-a5b6-66e54666c88e": {
		"id": "091b6583-0933-4d17-a5b6-66e54666c88e",
		"meta": map[string]interface{}{
			"version":      6,
			"created":      "2014-05-22T22:36:36.941Z",
			"lastModified": "2014-06-25T23:10:03.845Z",
		},
		"userName": "admin",
		"name": map[string]string{
			"familyName": "Admin",
			"givenName":  "Mister",
		},
		"emails": []map[string]string{
			{"value": "why-email@example.com"},
		},
		"groups": []map[string]string{
			{
				"value":   "e7f74565-4c7e-44ba-b068-b16072cbf08f",
				"display": "clients.read",
				"type":    "DIRECT",
			},
		},
		"approvals": []string{},
		"active":    true,
		"verified":  false,
		"schemas":   []string{"urn:scim:schemas:core:1.0"},
	},
	"943e6076-b1a5-4404-811b-a1ee9253bf56": {
		"id": "943e6076-b1a5-4404-811b-a1ee9253bf56",
		"meta": map[string]interface{}{
			"version":      6,
			"created":      "2014-05-22T22:36:36.941Z",
			"lastModified": "2014-06-25T23:10:03.845Z",
		},
		"userName": "some-user",
		"name": map[string]string{
			"familyName": "Some",
			"givenName":  "User",
		},
		"emails": []map[string]string{
			{"value": "slayer@example.com"},
		},
		"groups": []map[string]string{
			{
				"value":   "e7f74565-4c7e-44ba-b068-b16072cbf08f",
				"display": "clients.read",
				"type":    "DIRECT",
			},
		},
		"approvals": []string{},
		"active":    true,
		"verified":  false,
		"schemas":   []string{"urn:scim:schemas:core:1.0"},
	},
}

var allUsersResponse = []map[string]interface{}{
	UAAUsers["091b6583-0933-4d17-a5b6-66e54666c88e"],
	UAAUsers["943e6076-b1a5-4404-811b-a1ee9253bf56"],
}
