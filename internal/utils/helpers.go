package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// GenerateOrderNumber generates a unique order number
func GenerateOrderNumber() string {
	timestamp := time.Now().Format("20060102150405")
	randomBytes := make([]byte, 4)
	rand.Read(randomBytes)
	randomStr := hex.EncodeToString(randomBytes)
	return fmt.Sprintf("ORD-%s-%s", timestamp, strings.ToUpper(randomStr))
}

// SlugifyString converts a string to a URL-friendly slug
func SlugifyString(s string) string {
	// Convert to lowercase
	s = strings.ToLower(s)

	// Replace spaces and underscores with hyphens
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")

	// Remove all characters that are not alphanumeric or hyphens
	reg := regexp.MustCompile("[^a-z0-9-]+")
	s = reg.ReplaceAllString(s, "")

	// Remove duplicate hyphens
	reg = regexp.MustCompile("-+")
	s = reg.ReplaceAllString(s, "-")

	// Trim hyphens from beginning and end
	s = strings.Trim(s, "-")

	return s
}

// CalculateMaintenanceUntil calculates the maintenance expiry date
func CalculateMaintenanceUntil(months int) time.Time {
	return time.Now().AddDate(0, months, 0)
}

// IsMaintenanceExpired checks if maintenance has expired
func IsMaintenanceExpired(maintenanceUntil time.Time) bool {
	return time.Now().After(maintenanceUntil)
}

// DaysUntilExpiry calculates days until maintenance expires
func DaysUntilExpiry(maintenanceUntil time.Time) int {
	duration := time.Until(maintenanceUntil)
	return int(duration.Hours() / 24)
}

// TruncateString truncates a string to specified length
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength] + "..."
}
