package myrasec

import "testing"

func TestNew(t *testing.T) {
	key := "abc123"
	secret := "123abc"

	api, err := New(key, secret)
	if err != nil {
		t.Error("Unexpected error")
	}

	if api.key != key {
		t.Errorf("Expected key to be [%s] but got [%s]\n", key, api.key)
	}

	if api.secret != secret {
		t.Errorf("Expected key to be [%s] but got [%s]\n", key, api.key)
	}
}

func TestNewWithEmptyKey(t *testing.T) {
	key := ""
	secret := "123abc"

	_, err := New(key, secret)
	if err == nil {
		t.Error("Passing an empty key should fail")
	}

	if err.Error() != "Missing API credentials" {
		t.Errorf("Expected error message to be [%s] but got [%s]", "Missing API credentials", err.Error())
	}
}

func TestNewWithEmptySecret(t *testing.T) {
	key := "abc123"
	secret := ""

	_, err := New(key, secret)
	if err == nil {
		t.Error("Passing an empty secret should fail")
	}

	if err.Error() != "Missing API credentials" {
		t.Errorf("Expected error message to be [%s] but got [%s]", "Missing API credentials", err.Error())
	}
}

func TestNewWithEmptyParams(t *testing.T) {
	key := ""
	secret := ""

	_, err := New(key, secret)
	if err == nil {
		t.Error("Passing an empty key/secret should fail")
	}

	if err.Error() != "Missing API credentials" {
		t.Errorf("Expected error message to be [%s] but got [%s]", "Missing API credentials", err.Error())
	}
}

func TestSetUserAgent(t *testing.T) {
	key := "abc123"
	secret := "123abc"

	api, err := New(key, secret)
	if err != nil {
		t.Error("Unexpected error.")
	}

	if api.UserAgent != DefaultAPIUserAgent {
		t.Errorf("Expected default UserAgent to be [%s] but got [%s]", DefaultAPIUserAgent, api.UserAgent)
	}

	api.SetUserAgent("Testing")

	if api.UserAgent != "Testing" {
		t.Errorf("Expected UserAgent to be [%s] but got [%s]\n", "Testing", api.UserAgent)
	}
}

func TestSetLanguage(t *testing.T) {
	key := "abc123"
	secret := "123abc"

	api, err := New(key, secret)
	if err != nil {
		t.Error("Unexpected error.")
	}

	if api.Language != DefaultAPILanguage {
		t.Errorf("Expected default language to be [%s] but got [%s]\n", DefaultAPILanguage, api.Language)
	}

	err = api.SetLanguage("de")
	if err != nil {
		t.Errorf("Expected [%s] to be a valid language [%s]\n", "de", err.Error())
	}

	if api.Language != "de" {
		t.Errorf("Expected language to be [%s] but got [%s]\n", "de", api.Language)
	}

	err = api.SetLanguage("en")
	if err != nil {
		t.Errorf("Expected [%s] to be a valid language [%s]\n", "en", err.Error())
	}

	if api.Language != "en" {
		t.Errorf("Expected language to be [%s] but got [%s]\n", "en", api.Language)
	}

	err = api.SetLanguage("fr")
	if err == nil {
		t.Errorf("Expected [%s] to be invalid as a language setting\n", "fr")
	}
}
