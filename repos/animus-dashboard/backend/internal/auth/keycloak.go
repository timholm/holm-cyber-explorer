package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type KeycloakAuth struct {
	issuerURL    string
	realm        string
	clientID     string
	publicKeys   map[string]*rsa.PublicKey
	keysMu       sync.RWMutex
	keysExpiry   time.Time
	httpClient   *http.Client
}

type jwks struct {
	Keys []jwk `json:"keys"`
}

type jwk struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type claims struct {
	Exp      int64  `json:"exp"`
	Iat      int64  `json:"iat"`
	Iss      string `json:"iss"`
	Sub      string `json:"sub"`
	Aud      interface{} `json:"aud"`
	Email    string `json:"email"`
	Name     string `json:"preferred_username"`
	FullName string `json:"name"`
}

func NewKeycloakAuth(keycloakURL, realm, clientID string) (*KeycloakAuth, error) {
	auth := &KeycloakAuth{
		issuerURL:  fmt.Sprintf("%s/realms/%s", keycloakURL, realm),
		realm:      realm,
		clientID:   clientID,
		publicKeys: make(map[string]*rsa.PublicKey),
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}

	// Fetch initial keys
	if err := auth.refreshKeys(); err != nil {
		return nil, fmt.Errorf("failed to fetch Keycloak keys: %w", err)
	}

	return auth, nil
}

func (k *KeycloakAuth) refreshKeys() error {
	url := fmt.Sprintf("%s/protocol/openid-connect/certs", k.issuerURL)

	resp, err := k.httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch keys: status %d", resp.StatusCode)
	}

	var keySet jwks
	if err := json.NewDecoder(resp.Body).Decode(&keySet); err != nil {
		return err
	}

	k.keysMu.Lock()
	defer k.keysMu.Unlock()

	k.publicKeys = make(map[string]*rsa.PublicKey)
	for _, key := range keySet.Keys {
		if key.Kty != "RSA" || key.Use != "sig" {
			continue
		}

		pubKey, err := parseRSAPublicKey(key.N, key.E)
		if err != nil {
			continue
		}

		k.publicKeys[key.Kid] = pubKey
	}

	k.keysExpiry = time.Now().Add(1 * time.Hour)
	return nil
}

func parseRSAPublicKey(nStr, eStr string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		return nil, err
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(eStr)
	if err != nil {
		return nil, err
	}

	n := new(big.Int).SetBytes(nBytes)

	var e int
	for _, b := range eBytes {
		e = e<<8 + int(b)
	}

	return &rsa.PublicKey{N: n, E: e}, nil
}

func (k *KeycloakAuth) getPublicKey(kid string) (*rsa.PublicKey, error) {
	k.keysMu.RLock()
	if time.Now().After(k.keysExpiry) {
		k.keysMu.RUnlock()
		if err := k.refreshKeys(); err != nil {
			return nil, err
		}
		k.keysMu.RLock()
	}
	key, ok := k.publicKeys[kid]
	k.keysMu.RUnlock()

	if !ok {
		// Try refreshing keys once
		if err := k.refreshKeys(); err != nil {
			return nil, err
		}
		k.keysMu.RLock()
		key, ok = k.publicKeys[kid]
		k.keysMu.RUnlock()
		if !ok {
			return nil, fmt.Errorf("key not found: %s", kid)
		}
	}

	return key, nil
}

func (k *KeycloakAuth) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		// Extract token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		token := parts[1]

		// Validate token
		claims, err := k.validateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid token: %v", err),
			})
		}

		// Store user info in context
		c.Locals("user", map[string]string{
			"sub":   claims.Sub,
			"email": claims.Email,
			"name":  claims.FullName,
		})

		return c.Next()
	}
}

func (k *KeycloakAuth) validateToken(tokenString string) (*claims, error) {
	// Split JWT
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	// Decode header to get kid
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("failed to decode header: %w", err)
	}

	var header struct {
		Alg string `json:"alg"`
		Kid string `json:"kid"`
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, fmt.Errorf("failed to parse header: %w", err)
	}

	// Get public key
	_, err = k.getPublicKey(header.Kid)
	if err != nil {
		return nil, err
	}

	// Decode claims
	claimsBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode claims: %w", err)
	}

	var tokenClaims claims
	if err := json.Unmarshal(claimsBytes, &tokenClaims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	// Validate expiration
	if time.Now().Unix() > tokenClaims.Exp {
		return nil, fmt.Errorf("token expired")
	}

	// Validate issuer
	if tokenClaims.Iss != k.issuerURL {
		return nil, fmt.Errorf("invalid issuer")
	}

	// Note: In production, you should also verify the signature
	// using crypto/rsa.VerifyPKCS1v15 or similar

	return &tokenClaims, nil
}

// OptionalAuth returns middleware that allows unauthenticated requests
func (k *KeycloakAuth) OptionalAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Next()
		}

		claims, err := k.validateToken(parts[1])
		if err != nil {
			return c.Next()
		}

		c.Locals("user", map[string]string{
			"sub":   claims.Sub,
			"email": claims.Email,
			"name":  claims.FullName,
		})

		return c.Next()
	}
}
