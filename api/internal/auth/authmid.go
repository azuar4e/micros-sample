package auth

import (
	"context"
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

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawToken, err := c.Cookie("AccessToken")

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No authorization cookie found"})
			return
		}

		if rawToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Empty authorization token"})
			return
		}

		ctx := oidc.InsecureIssuerURLContext(context.Background(), "http://localhost:8080/realms/mini-e-comerce")
		provider, err := oidc.NewProvider(ctx, "http://keycloak-svc:8080/realms/mini-e-comerce")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "could not reach oidc provider", "details": err.Error()})
			return
		}

		verifier := provider.Verifier(&oidc.Config{SkipClientIDCheck: true})
		idToken, err := verifier.Verify(c.Request.Context(), rawToken)
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
