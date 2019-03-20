package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mattermost/mattermost-server/plugin"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/mattermost/mattermost-server/plugin/plugintest/mock"
	"github.com/stretchr/testify/assert"
)

var (
	basicConfig = &configuration{
		OAuthClientID:       "client-id",
		OAuthClientSecret:   "client-secret",
		WebexAuthorizeURL:   "https://api.ciscospark.com/v1/authorize",
		WebexAccessTokenURL: "https://api.ciscospark.com/v1/access_token",
		EncryptionKey:       "KoTkgdxmgP2XnJHRKsf_Dce-IZwP3nuX",
	}

	baseConnectedRequest = httptest.NewRequest("GET", "/connected", nil)
)

func TestConnected_Fail_NotAuthorizedNoHeaderMMUserID(t *testing.T) {

	req := baseConnectedRequest

	api := &plugintest.API{}

	p := Plugin{}
	p.setConfiguration(basicConfig)
	p.SetAPI(api)

	w := httptest.NewRecorder()
	p.ServeHTTP(&plugin.Context{}, w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
}

func TestConnected_Fail_WebexNotConnected(t *testing.T) {

	req := baseConnectedRequest
	req.Header.Set("Mattermost-User-ID", "abcd1234")

	api := &plugintest.API{}

	api.Mock.On("KVGet", mock.Anything).Return(nil, nil)
	api.Mock.On("KVDelete", mock.Anything).Return(nil)

	p := Plugin{}
	p.setConfiguration(basicConfig)
	p.SetAPI(api)

	w := httptest.NewRecorder()
	p.ServeHTTP(&plugin.Context{}, w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
}

func TestConnected_Success_WebexConnected(t *testing.T) {

	req := baseConnectedRequest
	req.Header.Set("Mattermost-User-ID", validUserId)

	api := &plugintest.API{}

	user := validUser

	sess := WebexOAuthSession{
		UserID: validUserId,
		Token:  validToken,
	}

	expected := map[string]interface{}{
		"user":    user,
		"session": sess,
	}

	var err error

	sess.Token.AccessToken, err = encrypt([]byte(basicConfig.EncryptionKey), sess.Token.AccessToken)
	assert.NoError(t, err)
	sess.Token.RefreshToken, err = encrypt([]byte(basicConfig.EncryptionKey), sess.Token.RefreshToken)
	assert.NoError(t, err)

	sessData, _ := json.Marshal(sess)
	userData, _ := json.Marshal(user)

	api.Mock.On("KVGet", fmt.Sprintf("%s%s", WebexUserKey, user.ID)).Return(userData, nil)
	api.Mock.On("KVGet", fmt.Sprintf("%s%s", WebexOAuthSessionKey, user.ID)).Return(sessData, nil)

	api.Mock.On("KVDelete", mock.Anything).Return(nil)
	api.Mock.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()

	p := Plugin{}
	p.setConfiguration(basicConfig)
	p.SetAPI(api)

	w := httptest.NewRecorder()

	p.ServeHTTP(&plugin.Context{}, w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	r := w.Result()
	body, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)

	want, _ := json.Marshal(expected)
	assert.Equal(t, string(want), string(body))
}
