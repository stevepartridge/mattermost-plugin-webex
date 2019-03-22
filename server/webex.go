package main

import (
	"golang.org/x/oauth2"

	webexteams "github.com/jbogarin/go-cisco-webex-teams/sdk"
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
