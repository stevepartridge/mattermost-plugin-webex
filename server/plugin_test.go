package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"time"

	"github.com/mattermost/mattermost-server/model"
	"golang.org/x/oauth2"
)

var (
	validUserId  = "abcd1234"
	validUserId2 = "efgh5678"

	validToken = oauth2.Token{
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

	validUser2 = WebexUserInfo{
		ID:          validUserId2,
		FirstName:   "First2",
		LastName:    "Last2",
		DisplayName: "First2 Last2",
	}

	startMeetingRequest = StartMeetingRequest{
		ChannelId: fmt.Sprintf("%s__%s", validUserId, validUserId2),
	}

	validMeeting = WebexMeeting{
		ID:            model.NewId(),
		ChannelID:     startMeetingRequest.ChannelId,
		FromUserID:    validUserId,
		FromWebexUser: validUser,
		ToUserID:      validUserId2,
		ToWebexUser:   validUser2,
	}

	baseOAuthConnectRequest  = httptest.NewRequest("GET", "/oauth2/connect", nil)
	baseOAuthCallbackRequest = httptest.NewRequest("GET", "/oauth2/callback", nil)

	baseAPIV1MeetingRequest = httptest.NewRequest("POST", "/api/v1/meetings", makeRequestBody(startMeetingRequest))
)

func makeRequestBody(v interface{}) io.ReadCloser {
	data, _ := json.Marshal(v)
	return ioutil.NopCloser(bytes.NewBuffer(data))
}
