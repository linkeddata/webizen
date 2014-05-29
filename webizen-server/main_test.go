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

func TestDeiu(t *testing.T) {
	testflight.WithServer(handler, func(r *testflight.Requester) {
		response := r.Get("/?q=https://deiu.rww.io/profile/card")
		assert.Equal(t, response.StatusCode, 200)
		assert.Equal(t, response.Body, "{\n  \"https://deiu.rww.io/profile/card#me\": {\n    \"img\": [\n      \"http://deiu.rww.io/profile/avatar.jpg\"\n    ],\n    \"name\": [\n      \"Andrei Vlad Sambra\"\n    ]\n  }\n}\n")
	})
	testflight.WithServer(handler, func(r *testflight.Requester) {
		response := r.Get("/?q=Andrei")
		assert.Equal(t, response.StatusCode, 200)
		assert.Equal(t, string(response.Body), "{\n  \"https://deiu.rww.io/profile/card#me\": {\n    \"img\": [\n      \"http://deiu.rww.io/profile/avatar.jpg\"\n    ],\n    \"mbox\": [\n      \"andrei.sambra@gmail.com\"\n    ],\n    \"name\": [\n      \"Andrei Vlad Sambra\"\n    ]\n  }\n}\n")
	})
}
