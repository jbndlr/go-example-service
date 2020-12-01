package rest

import (
	"encoding/json"
	"jbndlr/example/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireAuthentication : Middleware checking JWT authentication info.
func RequireAuthentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Request.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"message": "Unauthorized",
				})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Bad Request",
			})
			return
		}

		auth.Validate(cookie.Value)
		ctx.Next()
	}
}

// Authenticate : Handle POST-body-based authentication.
func Authenticate(ctx *gin.Context) (*auth.Subject, error) {
	var creds auth.Credentials
	err := json.NewDecoder(ctx.Request.Body).Decode(&creds)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request",
		})
		return &auth.Subject{}, err
	}

	subject, token, err := auth.Authenticate(&creds)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return &auth.Subject{}, err
	}

	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:    "token",
		Value:   token.JWT,
		Expires: token.Expires,
	})

	return subject, nil
}
