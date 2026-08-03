package service

import "testing"

func TestInitialAdminCredentials(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{name: "valid", username: "admin", password: "strong-password-123"},
		{name: "missing username", password: "strong-password-123", wantErr: true},
		{name: "short username", username: "ad", password: "strong-password-123", wantErr: true},
		{name: "short password", username: "admin", password: "short", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("GEEKAI_ADMIN_USERNAME", tt.username)
			t.Setenv("GEEKAI_ADMIN_PASSWORD", tt.password)
			username, password, err := initialAdminCredentials()
			if (err != nil) != tt.wantErr {
				t.Fatalf("initialAdminCredentials() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && (username != tt.username || password != tt.password) {
				t.Fatalf("initialAdminCredentials() = %q, %q", username, password)
			}
		})
	}
}
