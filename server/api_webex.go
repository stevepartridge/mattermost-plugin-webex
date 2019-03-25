package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mattermost/mattermost-server/model"
)

// https://developer.webex.com/docs/api/getting-started
// https://github.com/webex/react-ciscospark

func (p *Plugin) handleGetWebexUser(w http.ResponseWriter, r *http.Request) {
	sessionUserID := r.Header.Get("Mattermost-User-ID")
	if sessionUserID == "" {
		JSONErrorResponse(w, ErrNotAuthorized, http.StatusUnauthorized)
		return
	}

	user, err := p.loadWebexUser(sessionUserID)
	if err != nil {
		JSONResponse(w, err, http.StatusInternalServerError)
		return
	}

	JSONResponse(w, user, http.StatusOK)
}

type StartMeetingRequest struct {
	ChannelId string `json:"channel_id"`
	Personal  bool   `json:"personal"`
	Topic     string `json:"topic"`
	MeetingId int    `json:"meeting_id"`
}

func (p *Plugin) handleStartMeeting(w http.ResponseWriter, r *http.Request) {
	sessionUserID := r.Header.Get("Mattermost-User-Id")
	if sessionUserID == "" {
		JSONErrorResponse(w, ErrNotAuthorized, http.StatusUnauthorized)
		return
	}

	var req StartMeetingRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	_, err := p.loadWebexSession(sessionUserID)
	switch {
	case
		err == ErrWebexSessionNotFound,
		err == ErrWebexUserNotFound:

		p.API.LogError("Error retrieving session", "err", err.Error())
		JSONErrorResponse(w, ErrWebexNotConnected, http.StatusUnauthorized)
		return

	case err != nil:
		JSONErrorResponse(w, err, http.StatusInternalServerError)
		return

	}

	user, err := p.loadWebexUser(sessionUserID)
	switch {
	case
		err == ErrWebexSessionNotFound,
		err == ErrWebexUserNotFound:
		JSONErrorResponse(w, ErrWebexNotConnected, http.StatusUnauthorized)
		return
	case err != nil:
		JSONErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	channel, appErr := p.API.GetChannel(req.ChannelId)
	if appErr != nil {
		JSONErrorResponse(w, appErr, http.StatusInternalServerError) // maybe bad request
		return
	}

	if channel.Type != model.CHANNEL_DIRECT {
		JSONErrorResponse(w,
			ErrReplacer(ErrCreateMeetingTypeNotSupported, model.CHANNEL_DIRECT),
			http.StatusBadRequest,
		)
		return
	}

	// Not sure how reliable this is...
	toUserID := ""
	parts := strings.Split(channel.Name, "__")
	if len(parts) > 0 {
		toUserID = parts[0]
		if toUserID == sessionUserID && len(parts) > 1 {
			toUserID = parts[1]
		}
	}

	if toUserID == "" {
		JSONErrorResponse(w, ErrCreateMeetingToUserIdNotFound, http.StatusBadRequest)
		return
	}

	if toUserID == sessionUserID {
		JSONErrorResponse(w, ErrCreateMeetingFailedOwnChannel, http.StatusBadRequest)
		return
	}

	isGuest := false

	toWebexUser, err := p.loadWebexUser(toUserID)
	switch {
	case err == ErrWebexUserNotFound:
		isGuest = true
	case err != nil:
		JSONErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	meeting := WebexMeeting{
		ID:            model.NewId(),
		ChannelID:     channel.Id,
		FromUserID:    sessionUserID,
		FromWebexUser: *user,
		ToUserID:      toUserID,
	}

	if toWebexUser != nil {
		meeting.ToWebexUser = *toWebexUser
	}

	if isGuest {
		toUser, err := p.API.GetUser(toUserID)
		if err != nil {
			JSONErrorResponse(w, err, http.StatusInternalServerError)
			return
		}
		meeting.GuestEmail = toUser.Email
	}

	siteURL := p.API.GetConfig().ServiceSettings.SiteURL

	meeting.URL = fmt.Sprintf("%s/plugins/%s/meetings/%s", *siteURL, manifest.Id, meeting.ID)

	post := &model.Post{
		UserId:    sessionUserID,
		ChannelId: channel.Id,
		Message:   fmt.Sprintf("Webex Meeting started: %s", meeting.URL),
		Type:      "custom_webex",
		Props: map[string]interface{}{
			"meeting_id":        meeting.ID,
			"meeting_link":      meeting.URL,
			"meeting_status":    "STARTED",
			"meeting_personal":  false,
			"meeting_topic":     "Webex Meeting",
			"from_webhook":      "true",
			"override_username": "Webex",
			"override_icon_url": fmt.Sprintf("%s/static/plugins/%s/images/webex-logo.png", *siteURL, manifest.Id),
		},
	}

	_, appErr = p.API.CreatePost(post)
	if appErr != nil {
		JSONErrorResponse(w, appErr, http.StatusInternalServerError)
		return
	}

	err = p.storeWebexMeeting(meeting)
	if err != nil {
		JSONErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	JSONResponse(w, map[string]interface{}{"meeting": meeting}, http.StatusOK)
}
