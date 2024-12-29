package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// SanitizeURL checks that the URL is valid, safe, and matches the expected format for Twitch.
func SanitizeURL(inputURL string) (string, error) {
	// Ensure the URL is not empty
	if inputURL == "" {
		return "", fmt.Errorf("URL cannot be empty")
	}

	// Parse the URL to ensure it's valid
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL format: %w", err)
	}

	// Check if the URL is from Twitch (https://www.twitch.tv/)
	if parsedURL.Scheme != "https" || parsedURL.Host != "www.twitch.tv" {
		return "", fmt.Errorf("invalid Twitch URL, must start with https://www.twitch.tv/")
	}

	// Sanitize the URL by escaping special shell characters
	re := regexp.MustCompile(`[^a-zA-Z0-9\-._~:/?#[\]@!$&'()*+,;%=]`)
	sanitizedURL := re.ReplaceAllString(parsedURL.String(), "")

	// Optionally, you could verify that the sanitized URL still contains the necessary parts
	if !strings.HasPrefix(sanitizedURL, "https://www.twitch.tv/") {
		return "", fmt.Errorf("invalid Twitch URL, must start with https://www.twitch.tv/")
	}

	return sanitizedURL, nil
}
