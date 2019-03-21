package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/mattermost/mattermost-server/plugin/plugintest/mock"
	"github.com/stretchr/testify/assert"
)

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// )

// func TestOAuth2(t *testing.T) {

// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")
// 		w.Write([]byte(`{
//   "access_token":"ZDI3MGEyYzQtNmFlNS00NDNhLWFlNzAtZGVjNjE0MGU1OGZmZWNmZDEwN2ItYTU3",
//   "expires_in":1209600,
//   "refresh_token":"MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTEyMzQ1Njc4",
//   "refresh_token_expires_in":7776000
// }`))
// 	}))
// 	defer ts.Close()

// }

// {
//  "access_token":"ZDI3MGEyYzQtNmFlNS00NDNhLWFlNzAtZGVjNjE0MGU1OGZmZWNmZDEwN2ItYTU3",
//  "expires_in":1209600, //seconds
//  "refresh_token":"MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTEyMzQ1Njc4",
//  "refresh_token_expires_in":7776000 //seconds
// }

func TestGetOAuthConfig_Success(t *testing.T) {

	api := &plugintest.API{}
	p := Plugin{}

	siteURL := "http://example.com"
	cfg := &model.Config{
		ServiceSettings: model.ServiceSettings{
			SiteURL: &siteURL,
		},
	}

	api.Mock.On("GetConfig").Return(cfg)

	p.setConfiguration(basicConfig)
	p.SetAPI(api)

	err := p.OnActivate()
	assert.NoError(t, err)

	conf := p.getOAuthConfig()

	assert.Equal(t, conf.ClientID, basicConfig.OAuthClientID)

}

func Test_HandleOAuthConnect_Success(t *testing.T) {

	pluginConfig := basicConfig

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("req", r.RequestURI)
	}))
	defer ts.Close()

	pluginConfig.WebexAuthorizeURL = fmt.Sprintf("%s/v1/authorize", ts.URL)
	pluginConfig.WebexAccessTokenURL = fmt.Sprintf("%s/v1/access_token", ts.URL)

	req := baseOAuthConnectRequest
	req.Header.Set("Mattermost-User-ID", validUserId)

	api := &plugintest.API{}
	p := Plugin{}

	siteURL := "http://example.com"
	cfg := &model.Config{
		ServiceSettings: model.ServiceSettings{
			SiteURL: &siteURL,
		},
	}

	api.Mock.On("GetConfig").Return(cfg)
	api.Mock.On("KVSetWithExpiry", mock.Anything, mock.Anything, int64(300)).Return(nil)

	p.setConfiguration(pluginConfig)
	p.SetAPI(api)

	w := httptest.NewRecorder()

	p.ServeHTTP(&plugin.Context{}, w, req)

	assert.Equal(t, 302, w.Result().StatusCode)

}

func Test_HandleOAuthCallback_Success(t *testing.T) {

	pluginConfig := basicConfig

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

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

	pluginConfig.WebexAuthorizeURL = fmt.Sprintf("%s/v1/authorize", ts.URL)
	pluginConfig.WebexAccessTokenURL = fmt.Sprintf("%s/v1/access_token", ts.URL)

	state := fmt.Sprintf("%s_%s", model.NewId()[0:15], validUserId)

	authCode := "YjAzYzgyNDYtZTE3YS00OWZkLTg2YTgtNDc3Zjg4YzFiZDlkNTRlN2FhMjMtYzUz"

	req := baseOAuthCallbackRequest
	values := req.URL.Query()
	values.Add("code", authCode)
	values.Add("state", state)

	req.URL.RawQuery = values.Encode()

	api := &plugintest.API{}
	p := Plugin{}

	siteURL := "http://example.com"
	cfg := &model.Config{
		ServiceSettings: model.ServiceSettings{
			SiteURL: &siteURL,
		},
	}

	api.Mock.On("GetConfig").Return(cfg)
	api.Mock.On("KVGet", state).Return([]byte(state), nil)
	api.Mock.On("KVDelete", state).Return(nil)

	sessKey := fmt.Sprintf("%s%s", WebexOAuthSessionKey, validUserId)
	userKey := fmt.Sprintf("%s%s", WebexUserKey, validUserId)

	sess := WebexOAuthSession{
		UserID: validUserId,
		Token:  validToken,
	}
	sessData, _ := json.Marshal(sess)

	user := WebexUserInfo{}
	userData, _ := json.Marshal(user)

	api.Mock.On("KVSet", sessKey, mock.Anything).Return(nil)
	api.Mock.On("KVSet", userKey, mock.Anything).Return(nil)

	api.Mock.On("KVGet", sessKey).Return(sessData, nil)
	api.Mock.On("KVGet", userKey).Return(userData, nil)
	api.Mock.On("PublishWebSocketEvent", "oauth_success", mock.Anything, mock.Anything).Return()

	api.Mock.On("KVGet", fmt.Sprintf("%s_redir", state)).Return([]byte{}, nil)

	api.Mock.On("LogInfo", mock.Anything, mock.Anything, mock.Anything).Return()
	api.Mock.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Return()
	// api.Mock.On("KVSetWithExpiry", mock.Anything, mock.Anything, int64(300)).Return(nil)

	p.setConfiguration(pluginConfig)
	p.SetAPI(api)

	w := httptest.NewRecorder()

	p.ServeHTTP(&plugin.Context{}, w, req)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

}
