package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
)

func sign(payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func TestVerifyGitHubSignature(t *testing.T) {
	payload := []byte(`{"action":"created"}`)

	tests := []struct {
		name   string
		header string
		secret string
		want   bool
	}{
		{"valid", sign(payload, "s3cret"), "s3cret", true},
		{"wrong secret", sign(payload, "other"), "s3cret", false},
		{"wrong payload", sign([]byte(`{"action":"deleted"}`), "s3cret"), "s3cret", false},
		{"missing prefix", strings.TrimPrefix(sign(payload, "s3cret"), "sha256="), "s3cret", false},
		{"empty header", "", "s3cret", false},
		{"invalid hex", "sha256=zz", "s3cret", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := verifyGitHubSignature(payload, tt.header, tt.secret); got != tt.want {
				t.Fatalf("verifyGitHubSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
