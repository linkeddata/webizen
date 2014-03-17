package main

import (
	"github.com/drewolson/testflight"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test404(t *testing.T) {
	testflight.WithServer(handler, func(r *testflight.Requester) {
		response := r.Get("/?q=404")
		assert.Equal(t, response.StatusCode, 404)
	})
}
