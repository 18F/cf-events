package helpers

import (
	"fmt"
	"os"
	"testing"
)

type tokenConfigTest struct {
	key   string
	value string
}

var tokenConfigTests = []tokenConfigTest{
	{
		key:   "accept",
		value: "application/json",
	},
	{
		key:   "Content-Type",
		value: "application/x-www-form-urlencoded",
	},
	{
		key:   "authorization",
		value: "Basic Y2Y6",
	},
}

func TestTokenRequestConfig(t *testing.T) {
	req := config_token_request()
	// Check request header
	for _, test := range tokenConfigTests {
		if req.Header.Get(test.key) != test.value {
			t.Error("HTTP header not constructed properly expected %s but got %s", test.value, req.Header.Get(test.key))
		}
	}
	// Check request url
	expected_url := fmt.Sprintf("https://uaa.%s/oauth/token", os.Getenv("API_URL"))
	if fmt.Sprint(req.URL) != expected_url {
		t.Error("Expected URL to be %s, but got %s", expected_url, req.URL)
	}
}
