package auth

import "testing"

func TestAPIKeyAuthenticationWithoutOAuthFile(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("SPLITWISE_API_KEY", "test-api-key")
	token, err := LoadToken()
	if err != nil || token != "test-api-key" {
		t.Fatalf("API key authentication failed: %v", err)
	}
	if !IsLoggedIn() {
		t.Fatal("API key should count as authenticated")
	}
	t.Setenv("SPLITWISE_API_KEY", "")
	if _, err := LoadToken(); err == nil {
		t.Fatal("missing key and OAuth file should fail")
	}
	if IsLoggedIn() {
		t.Fatal("missing credentials should not count as authenticated")
	}
}
