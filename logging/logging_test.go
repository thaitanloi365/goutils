package logging

import (
	"testing"
)

func TestLogging(t *testing.T) {
	var logging = New()

	type st struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}

	logging.Debugw("Testing struct", st{ID: "1234", Title: "this is a struct"})
}
