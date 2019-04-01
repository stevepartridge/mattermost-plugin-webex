// +build unit_test

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	webexteams "github.com/jbogarin/go-cisco-webex-teams/sdk"
)

func NewWebexClient(token string) (*webexteams.Client, error) {

	if token == "" {
		return nil, ErrWebexClientMissingToken
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("req", r.RequestURI)
		if r.RequestURI == "/v1/access_token" {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{
  "access_token":"` + validToken.AccessToken + `",
  "expires_in":1209600,
  "refresh_token":"` + validToken.RefreshToken + `",
  "refresh_token_expires_in":7776000
}`))
		}

		if r.RequestURI == "/people/me" {
			w.Header().Set("Content-Type", "application/json")

			webexPerson := webexteams.Person{}

			if r.Header.Get("Authorization") != "" {
				if r.Header.Get("Authorization") == fmt.Sprintf("Bearer %s", validToken.AccessToken) {
					webexPerson = webexteams.Person{
						ID:     validWebexUser.ID,
						Emails: validWebexUser.Emails,
					}
				}
			}

			data, _ := json.Marshal(webexPerson)
			w.Write(data)
		}

	}))

	c := webexteams.NewClient(nil)

	webexteams.RestyClient.SetHostURL(ts.URL)
	webexteams.RestyClient.SetAuthToken(token)

	return c, nil
}
