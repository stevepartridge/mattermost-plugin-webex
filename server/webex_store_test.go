package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/mattermost/mattermost-server/plugin/plugintest/mock"
)

func TestLoadWebexSession(t *testing.T) {
	api := &plugintest.API{}

	sess := WebexOAuthSession{
		UserID: validUserId,
		Token:  validToken,
	}

	expiry := time.Now().Add(time.Minute * 5).UTC()
	sess.Token.Expiry = expiry

	expected := sess

	sessData := makeSessionData(sess)
	api.Mock.On("KVGet", fmt.Sprintf("%s%s", WebexOAuthSessionKey, sess.UserID)).Return(sessData, nil)

	p := Plugin{}
	p.setConfiguration(basicConfig)
	p.SetAPI(api)
	err := p.OnActivate()
	assert.NoError(t, err)

	actual, err := p.loadWebexSession(sess.UserID)
	assert.NoError(t, err)

	assert.EqualValues(t, &expected, actual)
}

func TestStoreWebexSession(t *testing.T) {

	t.Run("Success setting Session", func(t *testing.T) {

		api := &plugintest.API{}

		sess := WebexOAuthSession{
			UserID: validUserId,
			Token:  validToken,
		}

		api.Mock.On("KVSet", fmt.Sprintf("%s%s", WebexOAuthSessionKey, sess.UserID), mock.Anything).Return(nil)

		p := Plugin{}
		p.setConfiguration(basicConfig)
		p.SetAPI(api)
		err := p.OnActivate()
		assert.NoError(t, err)

		err = p.storeWebexSession(sess)
		assert.NoError(t, err)
	})

	t.Run("Fail Missing user ID", func(t *testing.T) {

		api := &plugintest.API{}

		sess := WebexOAuthSession{
			Token: validToken,
		}

		p := Plugin{}
		p.setConfiguration(basicConfig)
		p.SetAPI(api)
		err := p.OnActivate()
		assert.NoError(t, err)

		err = p.storeWebexSession(sess)
		assert.EqualError(t, err, ErrUnableToSaveSessionMissingUserID.Error())
	})

	t.Run("Fail Missing Access Token", func(t *testing.T) {

		api := &plugintest.API{}

		sess := WebexOAuthSession{
			UserID: validUserId,
		}

		p := Plugin{}
		p.setConfiguration(basicConfig)
		p.SetAPI(api)
		err := p.OnActivate()
		assert.NoError(t, err)

		err = p.storeWebexSession(sess)
		assert.EqualError(t, err, ErrUnableToSaveSessionMissingAccessToken.Error())

	})
}

func TestLoadWebexUser(t *testing.T) {
	api := &plugintest.API{}

	user := validWebexUser

	expected := validWebexUser

	userData, _ := json.Marshal(user)

	api.Mock.On("KVGet", fmt.Sprintf("%s%s", WebexUserKey, validUserId)).Return(userData, nil)

	p := Plugin{}
	p.setConfiguration(basicConfig)
	p.SetAPI(api)
	err := p.OnActivate()
	assert.NoError(t, err)

	actual, err := p.loadWebexUser(validUserId)
	assert.NoError(t, err)

	assert.EqualValues(t, &expected, actual)
}

func TestStoreWebexUser(t *testing.T) {
	api := &plugintest.API{}

	t.Run("Fail With Missing User ID", func(t *testing.T) {
		p := Plugin{}
		p.setConfiguration(basicConfig)
		p.SetAPI(api)
		err := p.OnActivate()
		assert.NoError(t, err)

		err = p.storeWebexUser("", nil)
		assert.EqualError(t, err, ErrUnableToSaveWebexUserMissingUserID.Error())
	})

	t.Run("Fail With Nil User", func(t *testing.T) {
		p := Plugin{}
		p.setConfiguration(basicConfig)
		p.SetAPI(api)
		err := p.OnActivate()
		assert.NoError(t, err)

		err = p.storeWebexUser(validUserId, nil)
		assert.EqualError(t, err, ErrUnableToSaveWebexUserMissingUser.Error())
	})

	t.Run("Success", func(t *testing.T) {

		user := validWebexUser

		api.Mock.On("KVSet", fmt.Sprintf("%s%s", WebexUserKey, validUserId), mock.Anything).Return(nil)

		p := Plugin{}
		p.setConfiguration(basicConfig)
		p.SetAPI(api)
		err := p.OnActivate()
		assert.NoError(t, err)

		err = p.storeWebexUser(validUserId, &user)
		assert.NoError(t, err)
	})
}

func TestGetWebexUserInfo(t *testing.T) {

	t.Run("Fail with Webex User Not Found", func(t *testing.T) {
		api := &plugintest.API{}

		sess := WebexOAuthSession{
			UserID: validUserId,
			Token:  validToken,
		}

		sessData := makeSessionData(sess)
		api.Mock.On("KVGet", fmt.Sprintf("%s%s", WebexOAuthSessionKey, sess.UserID)).Return(sessData, nil)

		// webexUserInfo := validUser
		// data, _ := json.Marshal(webexUserInfo)

		api.Mock.On("KVGet", fmt.Sprintf("%s%s", WebexUserKey, validUserId)).Return([]byte{}, nil)
		api.Mock.On("KVSet", fmt.Sprintf("%s%s", WebexUserKey, validUserId), mock.Anything).Return(nil)

		p := Plugin{}
		p.setConfiguration(basicConfig)
		p.SetAPI(api)
		err := p.OnActivate()
		assert.NoError(t, err)

		user, err := p.getWebexUserInfo(validUser.Id)
		assert.NoError(t, err)

		assert.NotNil(t, user)
	})

	t.Run("Succeed with valid user found", func(t *testing.T) {
		api := &plugintest.API{}

		webexUserInfo := validWebexUser
		data, _ := json.Marshal(webexUserInfo)
		api.Mock.On("KVGet", fmt.Sprintf("%s%s", WebexUserKey, validUserId)).Return(data, nil)
		api.Mock.On("KVSet", fmt.Sprintf("%s%s", WebexUserKey, validUserId), mock.Anything).Return(nil)

		p := Plugin{}
		p.setConfiguration(basicConfig)
		p.SetAPI(api)
		err := p.OnActivate()
		assert.NoError(t, err)

		user, err := p.getWebexUserInfo(validUserId)
		assert.NoError(t, err)

		assert.EqualValues(t, webexUserInfo.ID, user.ID)
	})

}

func TestLoadMeeting(t *testing.T) {

	t.Run("Fail With Meeting Not Found", func(t *testing.T) {

		api := &plugintest.API{}

		api.Mock.On("KVGet", fmt.Sprintf("%s%s", WebexMeetingKey, validMeeting.ID)).Return([]byte{}, nil)

		p := Plugin{}
		p.setConfiguration(basicConfig)
		p.SetAPI(api)
		err := p.OnActivate()
		assert.NoError(t, err)

		meeting, err := p.loadWebexMeeting(validMeeting.ID)
		assert.Nil(t, meeting)
		assert.EqualError(t, err, ErrWebexMeetingNotFound.Error())
	})

	t.Run("Success", func(t *testing.T) {

		api := &plugintest.API{}

		meetingInfo := validMeeting
		data, _ := json.Marshal(meetingInfo)
		api.Mock.On("KVGet", fmt.Sprintf("%s%s", WebexMeetingKey, validMeeting.ID)).Return(data, nil)

		p := Plugin{}
		p.setConfiguration(basicConfig)
		p.SetAPI(api)
		err := p.OnActivate()
		assert.NoError(t, err)

		meeting, err := p.loadWebexMeeting(validMeeting.ID)
		assert.NoError(t, err)

		assert.EqualValues(t, validMeeting.ID, meeting.ID)

	})
}

func TestStoreMeeting(t *testing.T) {

	t.Run("Fail With Meeting ID Missing", func(t *testing.T) {

		api := &plugintest.API{}

		p := Plugin{}
		p.setConfiguration(basicConfig)
		p.SetAPI(api)
		err := p.OnActivate()
		assert.NoError(t, err)

		err = p.storeWebexMeeting(WebexMeeting{})
		assert.EqualError(t, err, ErrUnableToSaveWebexMeetingMissingID.Error())
	})

}

func makeSessionData(sess WebexOAuthSession) []byte {
	if sess.Token.Expiry.IsZero() {
		expiry := time.Now().Add(time.Minute * 5).UTC()
		sess.Token.Expiry = expiry
	}

	sess.Token.AccessToken, _ = encrypt([]byte(basicConfig.EncryptionKey), sess.Token.AccessToken)

	sess.Token.RefreshToken, _ = encrypt([]byte(basicConfig.EncryptionKey), sess.Token.RefreshToken)

	sessData, _ := json.Marshal(sess)

	return sessData
}
