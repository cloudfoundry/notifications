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

	"github.com/cloudfoundry-incubator/notifications/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type UAA struct {
	server *httptest.Server
}

func NewUAA() UAA {
	router := mux.NewRouter()
	router.HandleFunc("/oauth/token", UAAPostOAuthToken).Methods("POST")
	router.HandleFunc("/token_key", UAAGetTokenKey).Methods("GET")
	router.HandleFunc("/Users", UAAGetUsers).Methods("GET")
	router.HandleFunc("/{anything:.*}", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("UAA ROUTE REQUEST ---> %+v\n", req)
		w.WriteHeader(http.StatusTeapot)
	}))
	return UAA{
		server: httptest.NewUnstartedServer(router),
	}
}

func (s UAA) Boot() {
	s.server.Start()
	os.Setenv("UAA_HOST", s.server.URL)
}

func (s UAA) Close() {
	s.server.Close()
}

func ReadFile(filename string) string {
	env := config.NewEnvironment()
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
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	token.Claims["client_id"] = clientID

	switch req.Form.Get("grant_type") {
	case "client_credentials":
		token.Claims["scope"] = []string{"notifications.write", "emails.write", "notification_preferences.admin", "critical_notifications.write", "notification_templates.admin"}
	case "authorization_code":
		token.Claims["scope"] = []string{"notification_preferences.read", "notification_preferences.write"}
		token.Claims["user_id"] = strings.TrimSuffix(req.Form.Get("code"), "-code")
	}

	tokenString, err := token.SignedString([]byte(ReadFile("/acceptance/fixtures/private.pem")))
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
		"value": ReadFile("/acceptance/fixtures/public.pem"),
	})
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
})

var UAAGetUsers = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
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

var UAAUsers = map[string]interface{}{
	"user-111": map[string]interface{}{
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
	"user-123": map[string]interface{}{
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
	"user-456": map[string]interface{}{
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
	"user-789": map[string]interface{}{
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
}
