package youtube

import (
	"bufio"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

// LoadCookiesFromFile loads cookies from a Netscape format cookies.txt file
// This format is exported by browser extensions like "Get cookies.txt"
func LoadCookiesFromFile(filepath string) ([]*http.Cookie, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open cookies file: %w", err)
	}
	defer file.Close()

	var cookies []*http.Cookie
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Netscape format: domain, flag, path, secure, expiration, name, value
		// Tab-separated
		parts := strings.Split(line, "\t")
		if len(parts) < 7 {
			continue
		}

		domain := parts[0]
		path := parts[2]
		secure := strings.ToLower(parts[3]) == "true"
		name := parts[5]
		value := parts[6]

		// Parse expiration
		var expires time.Time
		if expInt, err := strconv.ParseInt(parts[4], 10, 64); err == nil && expInt > 0 {
			expires = time.Unix(expInt, 0)
		}

		cookie := &http.Cookie{
			Name:     name,
			Value:    value,
			Domain:   domain,
			Path:     path,
			Secure:   secure,
			HttpOnly: true,
		}

		if !expires.IsZero() {
			cookie.Expires = expires
		}

		cookies = append(cookies, cookie)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading cookies file: %w", err)
	}

	return cookies, nil
}

// LoadCookiesFromString loads cookies from a Netscape format string
func LoadCookiesFromString(content string) ([]*http.Cookie, error) {
	var cookies []*http.Cookie
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Netscape format: domain, flag, path, secure, expiration, name, value
		parts := strings.Split(line, "\t")
		if len(parts) < 7 {
			continue
		}

		domain := parts[0]
		path := parts[2]
		secure := strings.ToLower(parts[3]) == "true"
		name := parts[5]
		value := parts[6]

		// Parse expiration
		var expires time.Time
		if expInt, err := strconv.ParseInt(parts[4], 10, 64); err == nil && expInt > 0 {
			expires = time.Unix(expInt, 0)
		}

		cookie := &http.Cookie{
			Name:     name,
			Value:    value,
			Domain:   domain,
			Path:     path,
			Secure:   secure,
			HttpOnly: true,
		}

		if !expires.IsZero() {
			cookie.Expires = expires
		}

		cookies = append(cookies, cookie)
	}

	return cookies, nil
}

// CreateCookieJar creates an http.CookieJar with the provided cookies
func CreateCookieJar(cookies []*http.Cookie) (http.CookieJar, error) {
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	// Group cookies by domain and set them
	youtubeURL, _ := url.Parse("https://www.youtube.com")
	googleURL, _ := url.Parse("https://www.google.com")

	var ytCookies, gCookies []*http.Cookie
	for _, cookie := range cookies {
		if strings.Contains(cookie.Domain, "youtube.com") {
			ytCookies = append(ytCookies, cookie)
		} else if strings.Contains(cookie.Domain, "google.com") {
			gCookies = append(gCookies, cookie)
		}
	}

	if len(ytCookies) > 0 {
		jar.SetCookies(youtubeURL, ytCookies)
	}
	if len(gCookies) > 0 {
		jar.SetCookies(googleURL, gCookies)
	}

	return jar, nil
}

// GetCookieHeader returns the Cookie header value for YouTube requests
func GetCookieHeader(cookies []*http.Cookie) string {
	var parts []string
	for _, cookie := range cookies {
		if strings.Contains(cookie.Domain, "youtube.com") || strings.Contains(cookie.Domain, "google.com") {
			parts = append(parts, fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
		}
	}
	return strings.Join(parts, "; ")
}

// FilterYouTubeCookies returns only cookies relevant for YouTube
func FilterYouTubeCookies(cookies []*http.Cookie) []*http.Cookie {
	var filtered []*http.Cookie
	for _, cookie := range cookies {
		domain := strings.ToLower(cookie.Domain)
		if strings.Contains(domain, "youtube.com") || strings.Contains(domain, "google.com") {
			filtered = append(filtered, cookie)
		}
	}
	return filtered
}
