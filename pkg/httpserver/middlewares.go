package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	requestIDKey = "requestID"
)

func requestIDMiddleware(ctx *gin.Context) {
	ctx.Set(requestIDKey, uuid.NewString())
}

func getSecurityHeadersMiddleware(allowedHost string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.Host != allowedHost {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
			return
		}

		ctx.Header("X-Frame-Options", "DENY")
		ctx.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
		ctx.Header("X-XSS-Protection", "1; mode=block")
		ctx.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		ctx.Header("Referrer-Policy", "strict-origin")
		ctx.Header("X-Content-Type-Options", "nosniff")
		ctx.Next()
	}
}
