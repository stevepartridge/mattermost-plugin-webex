// +build unit_test

package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	webexteams "github.com/jbogarin/go-cisco-webex-teams/sdk"
	"gopkg.in/resty.v1"
)

// NewWebexClient is a helper to create a new webex sdk client
// references:
//   https://github.com/jbogarin/go-cisco-webex-teams
//   https://github.com/jbogarin/go-cisco-webex-teams/blob/master/examples/people/main.go#L15-L18
func NewWebexClient(token string) (*webexteams.Client, error) {
	client := resty.New()

	client.SetAuthToken(token)

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

	}))
	defer ts.Close()

	client.SetHostURL(ts.URL)

	return webexteams.NewClient(client), nil
}
