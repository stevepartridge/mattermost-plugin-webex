package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mattermost/mattermost-server/plugin"
)

func (p *Plugin) prepareRoutes() {

	p.mux = chi.NewRouter()

	// Connceted check
	p.mux.Get("/connected", p.handleConnected)

	// OAuth
	p.mux.Get("/oauth2/connect", p.handleOAuthConnect)
	p.mux.Get("/oauth2/callback", p.handleOAuthCallback)

	// Meeting
	p.mux.Get("/meetings/{meeting_id}", p.handleMeeting)

	// API
	p.mux.Post("/api/v1/meetings", p.handleStartMeeting)
	p.mux.Get("/api/v1/user", p.handleGetWebexUser)

}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	config := p.getConfiguration()
	if err := config.IsValid(); err != nil {
		http.Error(w, "This plugin is not configured.", http.StatusInternalServerError)
		return
	}

	if p.mux == nil {
		p.prepareRoutes()
	}

	p.mux.ServeHTTP(w, r)
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
