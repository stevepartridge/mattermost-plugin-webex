// +build !unit_test

package main

import (
	webexteams "github.com/jbogarin/go-cisco-webex-teams/sdk"
)

// NewWebexClient is a helper to create a new webex sdk client
// references:
//   https://github.com/jbogarin/go-cisco-webex-teams
//   https://github.com/jbogarin/go-cisco-webex-teams/blob/master/examples/people/main.go#L15-L18
func NewWebexClient(token string) (*webexteams.Client, error) {
	if token == "" {
		return nil, ErrWebexClientMissingToken
	}

	c := webexteams.NewClient(nil)
	webexteams.RestyClient.SetAuthToken(token)

	return c, nil
}
