package main

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProxyMatching(t *testing.T) {
	p, err := NewProxy(&Entry{
		Source:     "*.example.com",
		DestFolder: "/tmp/foo",
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	for _, tc := range []string{
		"https://foo.example.com",
		"https://foo.example.com",
		"https://bar.example.com",
	} {
		req := httptest.NewRequest("GET", tc, nil)
		assert.True(t, p.Match(req), "failed to match %s", tc)
	}

	for _, tc := range []string{
		"https://example.com",
		"https://bad.com",
		"https://foo.bad.com",
		"https://foo.bar.example.com",
	} {
		req := httptest.NewRequest("GET", tc, nil)
		assert.False(t, p.Match(req), "incorrectly matched %s", tc)
	}
}
