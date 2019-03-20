package main

import "errors"

var (
	ErrNotAuthorized            = errors.New("Not authorized")
	ErrMethodNotAllowed         = errors.New("Method not allowed")
	ErrAuthroizationCodeMissing = errors.New("Authorization code is missing")
	ErrOAuthInvalidState        = errors.New("Invalid state")
	ErrOAuthRetrievingState     = errors.New("Error retrieving stored state")

	ErrWebexSessionNotFound = errors.New("Webex session not found")
	ErrWebexSessionExpired  = errors.New("Webex session has expired")

	ErrUnableToSaveSessionMissingUserID      = errors.New("Unable to save session, missing user ID")
	ErrUnableToSaveSessionMissingAccessToken = errors.New("Unable to save session, missing access token")

	ErrWebexUserNotFound                  = errors.New("Webex user not found")
	ErrUnableToSaveWebexUserMissingUserID = errors.New("Unable to save webex user, missing user ID")

	ErrWebexNotConnected = errors.New("Webex account not connected")

	ErrUnableToSaveWebexMeetingMissingID = errors.New("Unable to save session, missing meeting ID")

	ErrCreateMeetingFailedOwnChannel = errors.New("Unable to create new meeting with self")
	ErrCreateMeetingToUserIdNotFound = errors.New("Unable to create new meeting User Not Found")

	ErrWebexMeetingNotFound = errors.New("Webex meeting not found")
)
