package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebexUserInfoFromPerson(t *testing.T) {

	t.Run("Fail With Nil Person", func(t *testing.T) {
		w := WebexUserInfo{}
		w.FromWebexPerson(nil)
		assert.EqualValues(t, w.ID, "")
	})
}
