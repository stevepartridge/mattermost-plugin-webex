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

func TestHandleMeeting(t *testing.T) {

	t.Run("Fail with unauthorized", func(t *testing.T) {

		api := &plugintest.API{}

		p := Plugin{}
		p.setConfiguration(basicConfig)
		p.SetAPI(api)
		err := p.OnActivate()
		assert.NoError(t, err)

		req := baseMeetingRequest

		w := httptest.NewRecorder()
		p.ServeHTTP(&plugin.Context{}, w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)

	})

	t.Run("Success", func(t *testing.T) {

		api := &plugintest.API{}

		sess := WebexOAuthSession{
			UserID: validUserId,
			Token:  validToken,
		}

		siteURL := "http://example.com"
		cfg := &model.Config{
			ServiceSettings: model.ServiceSettings{
				SiteURL: &siteURL,
			},
		}

		api.Mock.On("GetConfig").Return(cfg)
		api.Mock.On("LogInfo", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.Mock.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()

		sessData := makeSessionData(sess)
		api.Mock.On("KVGet", fmt.Sprintf("%s%s", WebexOAuthSessionKey, sess.UserID)).Return(sessData, nil)

		sess2 := WebexOAuthSession{
			UserID: validUserId2,
			Token:  validToken,
		}

		sessData2 := makeSessionData(sess2)
		api.Mock.On("KVGet", fmt.Sprintf("%s%s", WebexOAuthSessionKey, sess.UserID)).Return(sessData2, nil)

		webexUserInfo := validUser
		data, _ := json.Marshal(webexUserInfo)

		api.Mock.On("KVGet", fmt.Sprintf("%s%s", WebexUserKey, validUser.Id)).Return(data, nil)

		meetingData, _ := json.Marshal(validMeeting)
		api.Mock.On("KVGet", fmt.Sprintf("%s%s", WebexMeetingKey, validMeeting.ID)).Return(meetingData, nil)

		webexUserInfo2 := validUser2
		toUserData, _ := json.Marshal(webexUserInfo2)
		api.Mock.On("KVGet", fmt.Sprintf("%s%s", WebexUserKey, validUser2.Id)).Return(toUserData, nil)

		api.Mock.On("GetUser", validUser.Id).Return(&validUser, nil)

		api.Mock.On("GetUser", validUser2.Id).Return(&validUser2, nil)

		p := Plugin{}
		p.setConfiguration(basicConfig)
		p.SetAPI(api)
		err := p.OnActivate()
		assert.NoError(t, err)

		req := baseMeetingRequest
		req.Header.Set("Mattermost-User-ID", sess.UserID)

		w := httptest.NewRecorder()
		p.ServeHTTP(&plugin.Context{}, w, req)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	})
}
