package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mattermost/mattermost-server/plugin"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/stretchr/testify/assert"
)

func TestServeHTTP(t *testing.T) {

	t.Run("Fail With Invalid Config", func(t *testing.T) {
		api := &plugintest.API{}

		p := Plugin{}
		p.setConfiguration(&configuration{})
		p.SetAPI(api)
		err := p.OnActivate()
		assert.Error(t, err)

		w := httptest.NewRecorder()

		p.ServeHTTP(&plugin.Context{}, w, baseConnectedRequest)
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})
}
