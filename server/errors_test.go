package main

import (
	"fmt"
	"testing"
)

func TestErrReplacer_SingleString_Success(t *testing.T) {

	expected := fmt.Sprintf(ErrCreateMeetingTypeNotSupported.Error(), "TYPE")

	err := ErrReplacer(ErrCreateMeetingTypeNotSupported, "TYPE")

	if err.Error() != expected {
		t.Errorf("Wanted: %s Got: %s", expected, err.Error())
	}
}
