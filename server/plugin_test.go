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
	validUserId       = "gcan9z6isjdp7na51rzji4raic"
	validUserId2      = "d6u7mfhb4pnx8p4orc35471b8a"
	validWebexUserId  = "abcd1234"
	validWebexUserId2 = "efgh5678"

	validToken = oauth2.Token{
		AccessToken:  "ZDI3MGEyYzQtNmFlNS00NDNhLWFlNzAtZGVjNjE0MGU1OGZmZWNmZDEwN2ItYTU3",
		RefreshToken: "MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTEyMzQ1Njc4",
		Expiry:       time.Now().Add(time.Second * 1209600),
	}

	validUser = model.User{
		Id: validUserId,
	}

	validWebexUser = WebexUserInfo{
		ID:          validWebexUserId,
		FirstName:   "First",
		LastName:    "Last",
		DisplayName: "First Last",
		Emails:      []string{"user@company.com"},
	}

	validUser2 = model.User{
		Id: validUserId2,
	}

	validWebexUser2 = WebexUserInfo{
		ID:          validWebexUserId2,
		FirstName:   "First2",
		LastName:    "Last2",
		DisplayName: "First2 Last2",
		Emails:      []string{"user2@company.com"},
	}

	startMeetingRequest = StartMeetingRequest{
		ChannelId: fmt.Sprintf("%s__%s", validUserId, validUserId2),
	}

	validMeeting = WebexMeeting{
		ID:            model.NewId(),
		ChannelID:     startMeetingRequest.ChannelId,
		FromUserID:    validUserId,
		FromWebexUser: validWebexUser,
		ToUserID:      validUserId2,
		ToWebexUser:   validWebexUser2,
	}

	basicConfig = &configuration{
		OAuthClientID:       "client-id",
		OAuthClientSecret:   "client-secret",
		WebexAuthorizeURL:   "https://api.ciscospark.com/v1/authorize",
		WebexAccessTokenURL: "https://api.ciscospark.com/v1/access_token",
		EncryptionKey:       "KoTkgdxmgP2XnJHRKsf_Dce-IZwP3nuX",
	}

	baseConnectedRequest = httptest.NewRequest("GET", "/connected", nil)

	baseOAuthConnectRequest  = httptest.NewRequest("GET", "/oauth2/connect", nil)
	baseOAuthCallbackRequest = httptest.NewRequest("GET", "/oauth2/callback", nil)

	baseAPIV1MeetingRequest      = httptest.NewRequest("POST", "/api/v1/meetings", makeRequestBody(startMeetingRequest))
	baseAPIV1GetWebexUserRequest = httptest.NewRequest("GET", "/api/v1/user", nil)

	baseMeetingRequest = httptest.NewRequest("GET", fmt.Sprintf("/meetings/%s", validMeeting.ID), nil)
)

func makeRequestBody(v interface{}) io.ReadCloser {
	data, _ := json.Marshal(v)
	return ioutil.NopCloser(bytes.NewBuffer(data))
}
