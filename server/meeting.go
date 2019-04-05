package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	"github.com/go-chi/chi"
)

const meetingHTMLTemplate = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf8">

  <title>Webex Meeting with {{.WithName}}</title>
  <link rel="stylesheet" href="{{.SiteURL}}/static/plugins/webex/external/spark.css">
</head>
<body>

  <div style="width: 100%; height: 100%;"
    id="space"
    data-toggle="ciscospark-space"
    data-initial-activity="meet"
    data-access-token="{{.AccessToken}}"
    {{.DestinationAttributes}}
    />

  <script src="{{.SiteURL}}/static/plugins/webex/external/spark.js"></script>
</body>
</html>`

var (
	meetingHTML = template.Must(template.New("meeting").Parse(meetingHTMLTemplate))
)

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

func (p *Plugin) handleMeeting(w http.ResponseWriter, r *http.Request) {
	sessionUserID := r.Header.Get("Mattermost-User-Id")
	if sessionUserID == "" {
		JSONErrorResponse(w, ErrNotAuthorized, http.StatusUnauthorized)
		return
	}

	siteURL := p.API.GetConfig().ServiceSettings.SiteURL

	// Ensure the current user is connected to webex
	session, err := p.loadWebexSession(sessionUserID)
	if err != nil {
		if err == ErrWebexUserNotFound || err == ErrWebexSessionNotFound {
			// Set up a redirect so we can bring them back to the proper location
			redirectURL := fmt.Sprintf("%s%s", *siteURL, r.RequestURI)
			oauthConnectURL := fmt.Sprintf(
				"%s/plugins/%s/oauth2/connect?redirect_to=%s",
				*siteURL,
				manifest.Id,
				url.QueryEscape(redirectURL),
			)

			p.API.LogDebug("redirect url", "url", redirectURL)
			http.Redirect(w, r, oauthConnectURL, http.StatusTemporaryRedirect)
			return
		}
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	meetingID := chi.URLParam(r, "meeting_id")
	p.API.LogInfo("Load meeting", "meeting_id", meetingID)

	meeting, err := p.loadWebexMeeting(meetingID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var (
		meetingWithName           string
		destinationDataAttributes string

		isGuest = false
	)

	// If the 'ToUserID' is the session user and there's a guest email
	// that means this user is being called as a guest (with their MM email)
	if meeting.ToUserID == sessionUserID && meeting.GuestEmail != "" {
		isGuest = true
	}

	// If current user is a guest and we have the email
	if meeting.GuestEmail != "" && !isGuest {

		meetingWithName = meeting.GuestEmail

		destinationDataAttributes = `
    data-to-person-email="` + meeting.GuestEmail + `"
  `
	}

	// If we're not meeting with a guest let's get their info
	if meetingWithName == "" {

		withWebexUser := meeting.ToWebexUser
		meetingWithUserID := meeting.ToUserID

		// If we're the current user switch up the to/from
		if sessionUserID == meeting.ToUserID {

			withWebexUser = meeting.FromWebexUser
			meetingWithUserID = meeting.FromUserID

		}

		withUser, err := p.API.GetUser(meetingWithUserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		meetingWithName = withUser.GetFullName()
		if strings.TrimSpace(meetingWithName) == "" {
			meetingWithName = withWebexUser.DisplayName
		}

		destinationDataAttributes = `
    data-destination-id="` + withWebexUser.ID + `"
    data-destination-type="userId"
    `
	}

	w.Header().Set("Content-Type", "text/html")
	err = meetingHTML.Execute(w, struct {
		WithName              string
		SiteURL               string
		AccessToken           string
		DestinationAttributes string
	}{
		WithName:              meetingWithName,
		SiteURL:               *siteURL,
		AccessToken:           session.Token.AccessToken,
		DestinationAttributes: destinationDataAttributes,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
