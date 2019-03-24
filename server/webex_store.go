package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// storeWebexSession saves the webex session to the KV store
// AccessToken and RefreshToken are both encrypted with the key
// set in the config EncryptionKey generated field
func (p *Plugin) storeWebexSession(session WebexOAuthSession) error {

	if session.UserID == "" {
		return ErrUnableToSaveSessionMissingUserID
	}

	if session.Token.AccessToken == "" {
		return ErrUnableToSaveSessionMissingAccessToken
	}

	config := p.getConfiguration()

	accessToken, encryptErr := encrypt([]byte(config.EncryptionKey), session.Token.AccessToken)
	if encryptErr != nil {
		return encryptErr
	}

	session.Token.AccessToken = accessToken

	// Encrypt the refresh token if it's present
	if session.Token.RefreshToken != "" {

		refreshToken, err := encrypt([]byte(config.EncryptionKey), session.Token.RefreshToken)
		if err != nil {
			return err
		}

		session.Token.RefreshToken = refreshToken
	}

	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	if err := p.API.KVSet(fmt.Sprintf("%s%s", WebexOAuthSessionKey, session.UserID), data); err != nil {
		return err
	}

	return nil

}

// loadWebexSession retrieves the webex session if present from the KV store
func (p *Plugin) loadWebexSession(userID string) (*WebexOAuthSession, error) {

	key := fmt.Sprintf("%s%s", WebexOAuthSessionKey, userID)

	data, appErr := p.API.KVGet(key)
	if appErr != nil {
		p.API.LogError(
			"Error retrieving webex session",
			"user_id", userID,
			"key", key,
			"error", appErr.Error(),
		)
		return nil, appErr
	}

	if len(data) == 0 {
		return nil, ErrWebexSessionNotFound
	}

	var session WebexOAuthSession

	err := json.Unmarshal(data, &session)
	if err != nil {
		return nil, err
	}

	config := p.getConfiguration()

	if time.Now().UTC().After(session.Token.Expiry) {
		kvErr := p.API.KVDelete(key)
		if kvErr != nil {
			return nil, kvErr
		}
		return nil, ErrWebexSessionExpired
	}

	accessToken, err := decrypt([]byte(config.EncryptionKey), session.Token.AccessToken)
	if err != nil {
		return &session, err
	}

	session.Token.AccessToken = accessToken

	// Decrypt refresh token if it's present
	if session.Token.RefreshToken != "" {

		refreshToken, err := decrypt([]byte(config.EncryptionKey), session.Token.RefreshToken)
		if err != nil {
			return &session, err
		}

		session.Token.RefreshToken = refreshToken
	}

	return &session, nil
}

// getWebexUserInfo retrieves the webex user info from the KV store and
// if not present will request from the webex API and subsequently save
// it to the KV store if successful
func (p *Plugin) getWebexUserInfo(userID string) (*WebexUserInfo, error) {

	webexUser, loadErr := p.loadWebexUser(userID)

	switch {
	case loadErr == ErrWebexUserNotFound:

		session, err := p.loadWebexSession(userID)
		if err != nil {
			p.API.LogError("error looking up session", "error", err.Error())
			return webexUser, err
		}

		webex, err := NewWebexClient(session.Token.AccessToken)
		if err != nil {
			p.API.LogError("Error creating new webex client", "error", err.Error())
			return webexUser, err
		}

		person, _, err := webex.People.GetMe() // don't need resp so supressing it
		if err != nil {
			p.API.LogError("Error calling people.GetMe()", "error", err.Error())
			return webexUser, err
		}

		webexUser = &WebexUserInfo{}
		webexUser.FromWebexPerson(person)

		err = p.storeWebexUser(userID, webexUser)
		if err != nil {
			p.API.LogError("Error saving webex user", "error", err.Error())
			return webexUser, err
		}

	case (loadErr != nil):
		return nil, loadErr
	}

	return webexUser, nil

}

// loadWebexUser loads the webex user info from the KV store if present
func (p *Plugin) loadWebexUser(userID string) (*WebexUserInfo, error) {

	key := fmt.Sprintf("%s%s", WebexUserKey, userID)

	data, appErr := p.API.KVGet(key)
	if appErr != nil {
		p.API.LogError(
			"Error retrieving webex user info",
			"user_id", userID,
			"key", key,
			"error", appErr,
		)
		return nil, appErr
	}

	if len(data) == 0 {
		return nil, ErrWebexUserNotFound
	}

	var webexUser WebexUserInfo

	err := json.Unmarshal(data, &webexUser)
	if err != nil {
		return nil, err
	}

	return &webexUser, nil
}

// storeWebexUser saves the user info to the KV store
func (p *Plugin) storeWebexUser(userID string, user *WebexUserInfo) error {

	if userID == "" {
		return ErrUnableToSaveWebexUserMissingUserID
	}

	if user == nil {
		return ErrUnableToSaveWebexUserMissingUser
	}

	key := fmt.Sprintf("%s%s", WebexUserKey, userID)

	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	kverr := p.API.KVSet(key, data)
	if kverr != nil {
		p.API.LogError("Error saving webex user info", "error", err.Error())
		return kverr
	}

	return nil

}

// loadWebexUser loads the webex meeting info from the KV store if present
func (p *Plugin) loadWebexMeeting(meetingID string) (*WebexMeeting, error) {

	key := fmt.Sprintf("%s%s", WebexMeetingKey, meetingID)

	data, appErr := p.API.KVGet(key)
	if appErr != nil {
		p.API.LogError(
			"Error retrieving webex meeting info",
			"user_id", meetingID,
			"key", key,
			"error", appErr,
		)
		return nil, appErr
	}

	if len(data) == 0 {
		return nil, ErrWebexMeetingNotFound
	}

	var webexMeeting WebexMeeting

	err := json.Unmarshal(data, &webexMeeting)
	if err != nil {
		return nil, err
	}

	return &webexMeeting, nil
}

// storeWebexUser saves the meeting info to the KV store
func (p *Plugin) storeWebexMeeting(meeting WebexMeeting) error {

	if meeting.ID == "" {
		return ErrUnableToSaveWebexMeetingMissingID
	}

	key := fmt.Sprintf("%s%s", WebexMeetingKey, meeting.ID)

	data, err := json.Marshal(meeting)
	if err != nil {
		return err
	}

	kverr := p.API.KVSetWithExpiry(key, data, meetingExpirySeconds)
	if kverr != nil {
		p.API.LogError("Error saving webex meeting info", "error", kverr.Error())
		return kverr
	}

	return nil

}
