package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	Sub         string `json:"sub"`
	Email       string `json:"email"`
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}

var oidcVerifier *oidc.IDTokenVerifier

func InitAuth(ctx context.Context) error {
	ctx = oidc.InsecureIssuerURLContext(ctx, "http://localhost:8080/realms/mini-e-commerce")

	provider, err := oidc.NewProvider(ctx, "http://keycloak-svc:8080/realms/mini-e-commerce")
	if err != nil {
		return fmt.Errorf("could not init oidc provider: %w", err)
	}

	oidcVerifier = provider.Verifier(&oidc.Config{SkipClientIDCheck: true})

	return nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawToken, err := c.Cookie("AccessToken")
		if err != nil || rawToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No authorization cookie found"})
			return
		}

		idToken, err := oidcVerifier.Verify(c.Request.Context(), rawToken) // reutiliza el verifier ya creado
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token", "details": err.Error()})
			return
		}

		var claims Claims
		if err := idToken.Claims(&claims); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not parse claims"})
			return
		}

		c.Set("userID", claims.Sub)
		c.Set("userEmail", claims.Email)
		c.Set("roles", claims.RealmAccess.Roles)
		c.Next()
	}
}
