package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"golang.org/x/oauth2"
)

func (p *Plugin) getOAuthConfig() *oauth2.Config {
	config := p.getConfiguration()

	siteURL := p.API.GetConfig().ServiceSettings.SiteURL

	return &oauth2.Config{
		ClientID:     config.OAuthClientID,
		ClientSecret: config.OAuthClientSecret,
		Scopes:       []string{"spark:all"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.WebexAuthorizeURL,
			TokenURL: config.WebexAccessTokenURL,
		},
		RedirectURL: fmt.Sprintf("%s/plugins/%s/oauth2/callback", *siteURL, manifest.Id),
	}
}

// handleOAuthConnect handles the kick of the OAuth2 consent flow
func (p *Plugin) handleOAuthConnect(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Mattermost-User-ID")
	if userID == "" {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	config := p.getOAuthConfig()

	state := fmt.Sprintf("%s_%s", model.NewId()[0:15], userID)

	// Expire the value in case they never finish the flow
	err := p.API.KVSetWithExpiry(state, []byte(state), stateExpirationSeconds)
	if err != nil {
		p.API.LogError("Error saving to KV store (with expiry)", "error", err.Error(), "message", err.Message)
		http.Error(w, err.Message, http.StatusInternalServerError)
		return
	}

	redirectURL := r.URL.Query().Get("redirect_to")
	if redirectURL != "" {

		// Expire the value in case they never finish the flow
		// Have to use _redir because _redirect exceeds 50 char limit for key length
		err = p.API.KVSetWithExpiry(fmt.Sprintf("%s_redir", state), []byte(redirectURL), stateExpirationSeconds)
		if err != nil {
			p.API.LogError("Error saving redirect to KV store (with expiry)", "error", err.Error(), "message", err.Message)
			http.Error(w, err.Message, http.StatusInternalServerError)
			return
		}

	}

	url := config.AuthCodeURL(state, oauth2.AccessTypeOffline)

	http.Redirect(w, r, url, http.StatusFound)
}

// handleOAuthCallback takes care of the callback if all went well with the
// identity provider (webex cloud) and they've sent the user back to us
func (p *Plugin) handleOAuthCallback(w http.ResponseWriter, r *http.Request) {

	var (
		code   string
		state  string
		userID string
	)

	code = r.URL.Query().Get("code")
	if strings.TrimSpace(code) == "" {
		p.API.LogError(ErrAuthroizationCodeMissing.Error())
		http.Error(w, ErrAuthroizationCodeMissing.Error(), http.StatusBadRequest)
		return
	}

	state = r.URL.Query().Get("state")

	storedState, kverr := p.API.KVGet(state)
	if kverr != nil {
		p.API.LogError(ErrOAuthRetrievingState.Error(), "error", kverr.Error())
		http.Error(w, ErrOAuthRetrievingState.Error(), http.StatusInternalServerError)
		return
	}

	if string(storedState) != state {
		p.API.LogError(ErrOAuthInvalidState.Error(), "stored_state", string(storedState), "state", state)
		http.Error(w, ErrOAuthInvalidState.Error(), http.StatusBadRequest)
		return
	}

	stateParts := strings.Split(state, "_")
	if len(stateParts) > 0 {
		userID = stateParts[1]
	}

	// TODO: figure out how to test this scenario
	if userID == "" {
		p.API.LogError("Invalid state", "stored_state", string(storedState), "state", state)
		http.Error(w, "Unable to determine user ID", http.StatusInternalServerError)
		return
	}

	p.API.KVDelete(state)

	oauthConfig := p.getOAuthConfig()

	ctx := context.Background()
	tkn, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		p.API.LogError("Error with token exchange", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session := WebexOAuthSession{
		UserID: userID,
		Token:  *tkn,
	}

	err = p.storeWebexSession(session)
	if err != nil {
		p.API.LogError("Error saving session", "error", err.Error())
		http.Error(w, "Error saving session", http.StatusInternalServerError)
		return
	}

	userInfo, err := p.getWebexUserInfo(userID)
	if err != nil {
		p.API.LogError("Error retrieving webex user info", "error", err.Error())
		http.Error(w, "Error retrieving webex user info", http.StatusInternalServerError)
		return
	}

	p.API.LogInfo("Webex User Connected", "webex_user_id", userInfo.ID)

	// publish an event to notify client(s)
	p.API.PublishWebSocketEvent(
		"oauth_success",
		map[string]interface{}{
			"success": true,
		},
		&model.WebsocketBroadcast{},
	)

	// see if we had a redirect url in the state
	val, kverr := p.API.KVGet(fmt.Sprintf("%s_redir", state))
	if kverr != nil {
		p.API.LogError("Error state", "error", kverr.Error())
	}

	if kverr == nil && len(val) > 0 {
		redirectURL, err := url.QueryUnescape(string(val))
		if err == nil {
			http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
			return
		}
	}

	p.API.LogDebug("token", "token", tkn.AccessToken)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(`<html>
  <title>Webex Connection Success</title>
  <body>
  Connected
  <script>
    window.close();
  </script>
  </body>
</html>
`))
}
