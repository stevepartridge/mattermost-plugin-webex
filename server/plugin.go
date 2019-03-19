package main

import (
	"sync"

	"github.com/mattermost/mattermost-server/plugin"
)

const (
	WebexOAuthSessionKey = "webex_session_"
	WebexUserKey         = "webex_user_"
	WebexMeetingKey      = "webex_meeting_"
)

var (
	// stateExpirationSeconds is how long in seconds to keep the value in the KV store
	stateExpirationSeconds = int64(60 * 5) // 5 min

	// meetingExpirySecons is a clean up meetings from KV store
	meetingExpirySeconds = int64(60 * 60 * 24) // 1 day
)

type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
}

// See https://developers.mattermost.com/extend/plugins/server/reference/

func (p *Plugin) OnActivate() error {
	config := p.getConfiguration()
	if err := config.IsValid(); err != nil {
		return err
	}

	return nil
}
