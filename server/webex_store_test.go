package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"golang.org/x/oauth2"

	"github.com/stretchr/testify/assert"

	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/mattermost/mattermost-server/plugin/plugintest/mock"
)

var (
	validUserId = "abcd1234"
	validToken  = oauth2.Token{
		AccessToken:  "ZDI3MGEyYzQtNmFlNS00NDNhLWFlNzAtZGVjNjE0MGU1OGZmZWNmZDEwN2ItYTU3",
		RefreshToken: "MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTEyMzQ1Njc4",
		Expiry:       time.Now().Add(time.Second * 1209600),
	}

	validUser = WebexUserInfo{
		ID:          validUserId,
		FirstName:   "First",
		LastName:    "Last",
		DisplayName: "First Last",
	}
)

func TestLoadWebexSession_Success(t *testing.T) {

	api := &plugintest.API{}

	sess := WebexOAuthSession{
		UserID: validUserId,
		Token:  validToken,
	}

	expiry := time.Now().Add(time.Minute * 5).UTC()
	sess.Token.Expiry = expiry

	expected := sess

	var err error

	sess.Token.AccessToken, err = encrypt([]byte(basicConfig.EncryptionKey), sess.Token.AccessToken)
	assert.NoError(t, err)
	sess.Token.RefreshToken, err = encrypt([]byte(basicConfig.EncryptionKey), sess.Token.RefreshToken)
	assert.NoError(t, err)

	sessData, _ := json.Marshal(sess)

	api.Mock.On("KVGet", fmt.Sprintf("%s%s", WebexOAuthSessionKey, sess.UserID)).Return(sessData, nil)

	p := Plugin{}

	p.setConfiguration(basicConfig)

	p.SetAPI(api)

	err = p.OnActivate()
	assert.NoError(t, err)

	actual, err := p.loadWebexSession(sess.UserID)
	assert.NoError(t, err)

	assert.EqualValues(t, expected, actual)

}

func TestStoreWebexSession_Success(t *testing.T) {

	api := &plugintest.API{}

	sess := WebexOAuthSession{
		UserID: validUserId,
		Token:  validToken,
	}

	expiry := time.Now().Add(time.Minute * 5).UTC()
	sess.Token.Expiry = expiry

	var err error

	sess.Token.AccessToken, err = encrypt([]byte(basicConfig.EncryptionKey), sess.Token.AccessToken)
	assert.NoError(t, err)
	sess.Token.RefreshToken, err = encrypt([]byte(basicConfig.EncryptionKey), sess.Token.RefreshToken)
	assert.NoError(t, err)

	api.Mock.On("KVSet", fmt.Sprintf("%s%s", WebexOAuthSessionKey, sess.UserID), mock.Anything).Return(nil)

	p := Plugin{}

	p.setConfiguration(basicConfig)

	p.SetAPI(api)

	err = p.OnActivate()
	assert.NoError(t, err)

	err = p.storeWebexSession(sess)
	assert.NoError(t, err)

}

func TestLoadWebexUser_Success(t *testing.T) {

	api := &plugintest.API{}

	user := validUser

	expected := validUser

	userData, _ := json.Marshal(user)

	api.Mock.On("KVGet", fmt.Sprintf("%s%s", WebexUserKey, user.ID)).Return(userData, nil)

	p := Plugin{}

	p.setConfiguration(basicConfig)

	p.SetAPI(api)

	err := p.OnActivate()
	assert.NoError(t, err)

	actual, err := p.loadWebexUser(user.ID)
	assert.NoError(t, err)

	assert.EqualValues(t, expected, actual)

}

func TestStoreWebexUser_Success(t *testing.T) {

	api := &plugintest.API{}

	user := validUser

	api.Mock.On("KVSet", fmt.Sprintf("%s%s", WebexUserKey, user.ID), mock.Anything).Return(nil)

	p := Plugin{}

	p.setConfiguration(basicConfig)

	p.SetAPI(api)

	err := p.OnActivate()
	assert.NoError(t, err)

	err = p.storeWebexUser(user.ID, user)
	assert.NoError(t, err)

}
