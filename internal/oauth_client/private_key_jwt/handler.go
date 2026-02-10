// Package main. handler provides http handlers for the private_key_jwt method for OAuth 2.0 authentication with Casdoor.
// It provides simple handlers such as /login, /callback, /logout, /dashboard just for testing.
package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/golang-jwt/jwt/v4"
)

// loginHandler simply gets the sign-in url of casdoor and redirects to it.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Generate authorization URL -> e.g.,
	// http://localhost:8000/login/oauth/authorize?
	// 			client_id=some_client_id
	//			&response_type=code
	//			&redirect_uri=http://localhost:8080/callback
	//			&scope=read
	//			&state=app-built-in
	urll := casdoorsdkClient.GetSigninUrl(redirectUrl)
	http.Redirect(w, r, urll, http.StatusFound)
}

// In your callbackHandler, replace the token retrieval line:
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	token, err := getOAuthTokenWithAssertion(code, state)
	if err != nil {
		http.Error(w, "Failed to get token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse Access token to get claims | Get user information - Validated JWT token
	claims := &casdoorsdk.Claims{}
	if _, _, err = jwt.NewParser().ParseUnverified(token.AccessToken, claims); err != nil {
		http.Error(w, "Failed to parse token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	claims.AccessToken = token.AccessToken

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     sessionKey,
		Value:    claims.Subject,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   3600,
	})

	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

// -------------------------------------------- Helper Functions --------------------------------------------

// getOAuthTokenWithAssertion exchanges the authorization code for tokens using private_key_jwt.
func getOAuthTokenWithAssertion(code, state string) (*casdoorsdk.Token, error) {
	// 1. Create the JWT assertion
	audience := strings.ReplaceAll(baseCasdoorUrl, "http://", "https://") // casdoor expects https token endpoint as per object/client_auth.go#AuthenticateClientByAssertion
	clientAssertion, err := createClientAssertionJWT(clientId, audience)
	if err != nil {
		return nil, fmt.Errorf("failed to create client assertion: %w", err)
	}

	// 2. Prepare the token request (as application/x-www-form-urlencoded)
	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("code", code)
	// client_id is NOT sent as a separate parameter when using client_assertion
	formData.Set("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	formData.Set("client_assertion", clientAssertion)
	// If using PKCE, also set the 'code_verifier'
	// formData.Set("code_verifier", yourCodeVerifier)

	// 3. Make the HTTP POST request
	tokenEndpoint := fmt.Sprintf("%s/api/login/oauth/access_token", baseCasdoorUrl)
	req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request failed: %s, body: %s", resp.Status, body)
	}

	// 4. Parse the response
	var tokenResp casdoorsdk.Token
	if err = json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}
	return &tokenResp, nil
}

// createClientAssertionJWT generates a JWT signed with the private key.
func createClientAssertionJWT(clientID string, audience string) (string, error) {
	const defaultExpiry = 5 * time.Minute // Max 5 min as per RFC
	now := time.Now()

	claims := jwt.RegisteredClaims{
		Issuer:    clientID,
		Subject:   clientID,
		Audience:  jwt.ClaimStrings{audience},
		ExpiresAt: jwt.NewNumericDate(now.Add(defaultExpiry)),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        fmt.Sprintf("jti-%d", now.UnixNano()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}
