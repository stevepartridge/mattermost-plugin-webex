package main

import (
	"golang.org/x/oauth2"

	webexteams "github.com/jbogarin/go-cisco-webex-teams/sdk"
	"gopkg.in/resty.v1"
)

type WebexUserInfo struct {
	ID          string   `json:"id"`
	Type        string   `json:"type"`
	Emails      []string `json:"emails"`
	DisplayName string   `json:"display_name"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	Nickname    string   `json:"nickname"`
	Avatar      string   `json:"avatar"`
	Timezone    string   `json:"timezone"`

	Roles []string `json:"roles"`
	OrgID string   `json:"org_id"`
}

// FromWebexPerson is a helper to build out a WebexUserInfo from a webex object
func (self *WebexUserInfo) FromWebexPerson(person *webexteams.Person) {
	self.ID = person.ID
	self.Emails = person.Emails
	self.DisplayName = person.DisplayName
	self.FirstName = person.FirstName
	self.LastName = person.LastName
	self.Nickname = person.NickName
	self.Avatar = person.Avatar
	self.Timezone = person.TimeZone
	self.Type = person.PersonType

	self.OrgID = person.OrgID
	self.Roles = person.Roles
}

// WebexOAuthSession
type WebexOAuthSession struct {
	UserID string       `json:"user_id"`
	Token  oauth2.Token `json:"token"`
}

// WebexMeeting
type WebexMeeting struct {
	ID            string        `json:"id"`
	FromUserID    string        `json:"from_user_id"`
	FromWebexUser WebexUserInfo `json:"from_webex_user"`
	ChannelID     string        `json:"channel_id"`
	ToUserID      string        `json:"to_user_id"`
	ToWebexUser   WebexUserInfo `json:"to_webex_user"`
	GuestEmail    string        `json:"guest_email"`
	URL           string        `json:"meeting_url"`
}

// NewWebexClient is a helper to create a new webex sdk client
// references:
//   https://github.com/jbogarin/go-cisco-webex-teams
//   https://github.com/jbogarin/go-cisco-webex-teams/blob/master/examples/people/main.go#L15-L18
func NewWebexClient(token string) (*webexteams.Client, error) {

	client := resty.New()

	client.SetAuthToken(token)

	return webexteams.NewClient(client), nil
}
