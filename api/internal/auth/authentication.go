package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type AuthURL struct {
	URL          string
	State        string
	CodeVerifier string
}

func (u *AuthURL) String() string {
	return u.URL
}

var pendingAuths = struct {
	sync.Mutex
	m map[string]*AuthURL
}{m: make(map[string]*AuthURL)}

func oauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     "oidc-client",
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:8080/realms/mini-e-comerce/protocol/openid-connect/auth",
			TokenURL: "http://keycloak-svc:8080/realms/mini-e-comerce/protocol/openid-connect/token",
		},
		RedirectURL: "https://localhost:8443/api/v1/callback",
		Scopes:      []string{"openid", "profile", "email"},
	}
}

func Login(c *gin.Context) {
	cfg := oauthConfig()

	auth, err := AuthorizationURL(cfg)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	pendingAuths.Lock()
	pendingAuths.m[auth.State] = auth
	pendingAuths.Unlock()

	c.Redirect(http.StatusTemporaryRedirect, auth.URL)
}

func Callback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Code"})
		return
	}

	pendingAuths.Lock()
	auth, ok := pendingAuths.m[state]
	if ok {
		delete(pendingAuths.m, state)
	}
	pendingAuths.Unlock()

	if !ok || subtle.ConstantTimeCompare([]byte(auth.State), []byte(state)) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid State"})
		return
	}

	cfg := oauthConfig()

	token, err := cfg.Exchange(
		c.Request.Context(),
		code,
		oauth2.SetAuthURLParam("code_verifier", auth.CodeVerifier),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.SetCookie("AccessToken", token.AccessToken, 3600, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/api/v1/me")
}

func randomBytesInHex(count int) (string, error) {
	buf := make([]byte, count)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		return "", fmt.Errorf("Could not generate %d random bytes: %v", count, err)
	}

	return hex.EncodeToString(buf), nil
}

func AuthorizationURL(config *oauth2.Config) (*AuthURL, error) {
	codeVerifier, verifierErr := randomBytesInHex(32) // 64 character string here
	if verifierErr != nil {
		return nil, fmt.Errorf("Could not create a code verifier: %v", verifierErr)
	}
	sha2 := sha256.New()
	io.WriteString(sha2, codeVerifier)
	codeChallenge := base64.RawURLEncoding.EncodeToString(sha2.Sum(nil))

	state, stateErr := randomBytesInHex(24)
	if stateErr != nil {
		return nil, fmt.Errorf("Could not generate random state: %v", stateErr)
	}

	authUrl := config.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
	)

	return &AuthURL{
		URL:          authUrl,
		State:        state,
		CodeVerifier: codeVerifier,
	}, nil
}
