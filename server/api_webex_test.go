package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/mattermost/mattermost-server/plugin/plugintest/mock"
	"github.com/stretchr/testify/assert"
)

func TestConnected(t *testing.T) {
	req := baseConnectedRequest

	t.Run("Fail with not authorized no header Mattermost-User-ID", func(t *testing.T) {

		api := &plugintest.API{}

		p := Plugin{}
		p.setConfiguration(basicConfig)
		p.SetAPI(api)
		err := p.OnActivate()
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		p.ServeHTTP(&plugin.Context{}, w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})

	t.Run("Fail with Webex Not Connected", func(t *testing.T) {

		api := &plugintest.API{}

		p := Plugin{}
		p.setConfiguration(basicConfig)
		p.SetAPI(api)
		err := p.OnActivate()
		assert.NoError(t, err)

		req.Header.Set("Mattermost-User-ID", validUserId)

		api.Mock.On("KVGet", mock.Anything).Return(nil, nil)
		api.Mock.On("KVDelete", mock.Anything).Return(nil)

		w := httptest.NewRecorder()
		p.ServeHTTP(&plugin.Context{}, w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})

	t.Run("Success with Webex Connected", func(t *testing.T) {

		api := &plugintest.API{}

		p := Plugin{}
		p.setConfiguration(basicConfig)
		p.SetAPI(api)
		err := p.OnActivate()
		assert.NoError(t, err)

		user := validWebexUser

		sess := WebexOAuthSession{
			UserID: validUserId,
			Token:  validToken,
		}

		expected := map[string]interface{}{
			"user":    user,
			"session": sess,
		}

		sess.Token.AccessToken, err = encrypt([]byte(basicConfig.EncryptionKey), sess.Token.AccessToken)
		assert.NoError(t, err)
		sess.Token.RefreshToken, err = encrypt([]byte(basicConfig.EncryptionKey), sess.Token.RefreshToken)
		assert.NoError(t, err)

		sessData, _ := json.Marshal(sess)
		userData, _ := json.Marshal(user)

		api.Mock.On("KVGet", fmt.Sprintf("%s%s", WebexUserKey, validUserId)).Return(userData, nil)
		api.Mock.On("KVGet", fmt.Sprintf("%s%s", WebexOAuthSessionKey, validUserId)).Return(sessData, nil)

		api.Mock.On("KVDelete", mock.Anything).Return(nil)
		api.Mock.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()

		w := httptest.NewRecorder()

		req.Header.Set("Mattermost-User-ID", validUserId)

		p.ServeHTTP(&plugin.Context{}, w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		r := w.Result()
		body, err := ioutil.ReadAll(r.Body)
		assert.NoError(t, err)

		want, _ := json.Marshal(expected)
		assert.Equal(t, string(want), string(body))
	})

}

func TestStartMeeting(t *testing.T) {
	req := baseAPIV1MeetingRequest
	req.Header.Set("Mattermost-User-ID", validUserId)

	api := &plugintest.API{}

	p := Plugin{}
	p.setConfiguration(basicConfig)
	err := p.OnActivate()
	assert.NoError(t, err)

	sess := WebexOAuthSession{
		UserID: validUserId,
		Token:  validToken,
	}

	sessData := makeSessionData(sess)
	userData, _ := json.Marshal(validWebexUser)
	userData2, _ := json.Marshal(validWebexUser2)

	siteURL := "http://example.com"
	cfg := &model.Config{
		ServiceSettings: model.ServiceSettings{
			SiteURL: &siteURL,
		},
	}

	api.Mock.On("GetConfig").Return(cfg)

	api.Mock.On("KVGet", fmt.Sprintf("%s%s", WebexOAuthSessionKey, validUserId)).Return(sessData, nil)
	api.Mock.On("KVGet", fmt.Sprintf("%s%s", WebexUserKey, validUserId)).Return(userData, nil)
	api.Mock.On("KVGet", fmt.Sprintf("%s%s", WebexUserKey, validUserId2)).Return(userData2, nil)

	api.Mock.On("KVDelete", mock.Anything).Return(nil)
	api.Mock.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()

	channel := &model.Channel{
		Id:   startMeetingRequest.ChannelId,
		Name: startMeetingRequest.ChannelId,
		Type: model.CHANNEL_DIRECT,
	}

	api.Mock.On("GetChannel", startMeetingRequest.ChannelId).Return(channel, nil)

	api.Mock.On("CreatePost", mock.Anything).Return(nil, nil).Maybe()
	api.Mock.On("KVSetWithExpiry", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Maybe()

	p.SetAPI(api)

	w := httptest.NewRecorder()

	p.ServeHTTP(&plugin.Context{}, w, req)

	r := w.Result()
	_, err = ioutil.ReadAll(r.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

}

func TestAPIV1GetWebexUserInfo(t *testing.T) {

	t.Run("Fail with unauthorized", func(t *testing.T) {

		api := &plugintest.API{}

		p := Plugin{}
		p.setConfiguration(basicConfig)
		p.SetAPI(api)
		err := p.OnActivate()
		assert.NoError(t, err)

		req := baseAPIV1GetWebexUserRequest

		api.Mock.On("KVGet", mock.Anything).Return(nil, nil)

		w := httptest.NewRecorder()
		p.ServeHTTP(&plugin.Context{}, w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)

	})

	t.Run("Success with user info", func(t *testing.T) {

		api := &plugintest.API{}

		sess := WebexOAuthSession{
			UserID: validUserId,
			Token:  validToken,
		}

		sessData := makeSessionData(sess)
		api.Mock.On("KVGet", fmt.Sprintf("%s%s", WebexOAuthSessionKey, sess.UserID)).Return(sessData, nil)

		webexUserInfo := validWebexUser
		data, _ := json.Marshal(webexUserInfo)

		api.Mock.On("KVGet", fmt.Sprintf("%s%s", WebexUserKey, validUserId)).Return(data, nil)
		// api.Mock.On("KVSet", fmt.Sprintf("%s%s", WebexUserKey, validUser.ID), mock.Anything).Return(nil)

		p := Plugin{}
		p.setConfiguration(basicConfig)
		p.SetAPI(api)
		err := p.OnActivate()
		assert.NoError(t, err)

		req := baseAPIV1GetWebexUserRequest
		req.Header.Set("Mattermost-User-ID", sess.UserID)

		w := httptest.NewRecorder()
		p.ServeHTTP(&plugin.Context{}, w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	})
}
