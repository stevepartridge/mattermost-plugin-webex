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

func TestStartMeeting_Success(t *testing.T) {
	req := baseAPIV1MeetingRequest
	req.Header.Set("Mattermost-User-ID", validUserId)

	api := &plugintest.API{}

	p := Plugin{}
	p.setConfiguration(basicConfig)

	sess := WebexOAuthSession{
		UserID: validUserId,
		Token:  validToken,
	}

	sessData := makeSessionData(sess)
	userData, _ := json.Marshal(validUser)
	userData2, _ := json.Marshal(validUser2)

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
	}

	api.Mock.On("GetChannel", startMeetingRequest.ChannelId).Return(channel, nil)

	api.Mock.On("CreatePost", mock.Anything).Return(nil, nil).Maybe()
	api.Mock.On("KVSetWithExpiry", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Maybe()

	p.SetAPI(api)

	w := httptest.NewRecorder()

	p.ServeHTTP(&plugin.Context{}, w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)

	r := w.Result()
	_, err := ioutil.ReadAll(r.Body)
	assert.NoError(t, err)

	// want, _ := json.Marshal(expected)
	// assert.Equal(t, string(want), string(body))
}
