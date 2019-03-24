package main

import (
	"net/http"
	"strings"

	"github.com/mattermost/mattermost-server/plugin"
)

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	config := p.getConfiguration()
	if err := config.IsValid(); err != nil {
		http.Error(w, "This plugin is not configured.", http.StatusInternalServerError)
		return
	}

	path := r.URL.Path

	switch {
	// meeting route with id
	case strings.HasPrefix(path, "/meetings/"):
		p.handleMeeting(w, r)
		return
	}

	switch path {

	// connect check
	case "/connected":
		p.handleConnected(w, r)

	// oauth routes
	case "/oauth2/connect":
		p.handleOAuthConnect(w, r)
	case "/oauth2/callback":
		p.handleOAuthCallback(w, r)

	// api routes
	case "/api/v1/meetings":
		p.handleStartMeeting(w, r)
	case "/api/v1/user":
		p.handleGetWebexUser(w, r)
	default:
		http.NotFound(w, r)
	}

}

// handleConnected checks if the user has connected their webex account
func (p *Plugin) handleConnected(w http.ResponseWriter, r *http.Request) {

	requestorID := r.Header.Get("Mattermost-User-ID")
	if requestorID == "" {
		JSONErrorResponse(w, ErrNotAuthorized, http.StatusUnauthorized)
		return
	}

	user, err := p.loadWebexUser(requestorID)
	switch {
	case err == ErrWebexUserNotFound:
		JSONErrorResponse(w, ErrWebexNotConnected, http.StatusUnauthorized)
		return
	case err != nil:
		JSONErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	session, err := p.loadWebexSession(requestorID)
	if err != nil {
		JSONErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	JSONResponse(w, map[string]interface{}{
		"user":    user,
		"session": session,
	}, http.StatusOK)

}
