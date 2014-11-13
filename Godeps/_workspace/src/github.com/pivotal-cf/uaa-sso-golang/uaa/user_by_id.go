package uaa

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type UserByIDInterface interface {
	UserByID(string) (User, error)
}

type User struct {
	Username string
	ID       string
	Name     Name
	Emails   []string
	Active   bool
	Verified bool
}

type Name struct {
	FamilyName string
	GivenName  string
}

func UserByID(u UAA, id string) (User, error) {
	user := User{
		ID: id,
	}

	uri, err := url.Parse(fmt.Sprintf("%s/Users/%s", u.uaaURL, id))
	if err != nil {
		return user, err
	}

	host := uri.Scheme + "://" + uri.Host
	client := NewClient(host, u.VerifySSL).WithAuthorizationToken(u.AccessToken)
	code, body, err := client.MakeRequest("GET", uri.RequestURI(), nil)
	if err != nil {
		return user, err
	}

	if code > 399 {
		return user, NewFailure(code, body)
	}

	user, err = UserFromJSON(body)
	if err != nil {
		return user, err
	}

	return user, nil
}

func UserFromJSON(jsonBytes []byte) (User, error) {
	var parsed map[string]interface{}
	err := json.Unmarshal(jsonBytes, &parsed)
	if err != nil {
		return User{}, err
	}

	return UserFromResource(parsed)
}

func UserFromResource(resource map[string]interface{}) (User, error) {
	user := User{}

	userName, ok := resource["userName"].(string)
	if ok {
		user.Username = userName
	}

	id, ok := resource["id"].(string)
	if ok {
		user.ID = id
	}

	active, ok := resource["active"].(bool)
	if ok {
		user.Active = active
	}

	verified, ok := resource["verified"].(bool)
	if ok {
		user.Verified = verified
	}

	name, ok := resource["name"].(map[string]interface{})
	if ok {
		givenName, ok := name["givenName"].(string)
		if ok {
			user.Name.GivenName = givenName
		}

		familyName, ok := name["familyName"].(string)
		if ok {
			user.Name.FamilyName = familyName
		}
	}

	emailInterfaces, ok := resource["emails"].([]interface{})
	if ok {
		for _, emailInterface := range emailInterfaces {
			emailMap, ok := emailInterface.(map[string]interface{})
			if ok {
				email, ok := emailMap["value"].(string)
				if ok {
					user.Emails = append(user.Emails, email)
				}
			}
		}
	}
	return user, nil
}
