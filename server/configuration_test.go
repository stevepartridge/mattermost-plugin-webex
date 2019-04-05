package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfiguration(t *testing.T) {

	t.Run("Fail with Missing OAuth Client ID", func(t *testing.T) {

		cfg := *basicConfig
		cfg.OAuthClientID = ""

		p := Plugin{}
		p.setConfiguration(&cfg)
		err := p.OnActivate()
		assert.EqualError(t, err, ErrConfigInvalidMissingOAuthClientID.Error())

	})

	t.Run("Fail with Missing OAuth Client Secret", func(t *testing.T) {

		cfg := *basicConfig
		cfg.OAuthClientSecret = ""

		p := Plugin{}
		p.setConfiguration(&cfg)
		err := p.OnActivate()
		assert.EqualError(t, err, ErrConfigInvalidMissingOAuthClientSecret.Error())

	})

	t.Run("Success with Missing Authorize URL", func(t *testing.T) {

		cfg := *basicConfig
		cfg.WebexAuthorizeURL = ""

		p := Plugin{}
		p.setConfiguration(&cfg)
		err := p.OnActivate()
		assert.NoError(t, err)

	})

	t.Run("Success with Missing Access Token URL", func(t *testing.T) {

		cfg := *basicConfig
		cfg.WebexAccessTokenURL = ""

		p := Plugin{}
		p.setConfiguration(&cfg)
		err := p.OnActivate()
		assert.NoError(t, err)

	})

	t.Run("Fail with Missing Encryption Key", func(t *testing.T) {

		cfg := *basicConfig
		cfg.EncryptionKey = ""

		p := Plugin{}
		p.setConfiguration(&cfg)
		err := p.OnActivate()
		assert.EqualError(t, err, ErrConfigInvalidMissingEncryptionKey.Error())

	})
}
