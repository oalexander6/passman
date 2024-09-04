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
		// if ctx.Request.Header["Origin"][0] != allowedHost {
		// 	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
		// 	return
		// }

		ctx.Header("X-Frame-Options", "DENY")
		ctx.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
		ctx.Header("X-XSS-Protection", "1; mode=block")
		ctx.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		ctx.Header("Referrer-Policy", "strict-origin")
		ctx.Header("X-Content-Type-Options", "nosniff")

		ctx.Next()
	}
}

func csrfHeaderMiddleware(ctx *gin.Context) {
	if ctx.Request.Method != "OPTIONS" {
		if val := ctx.Request.Header["X-Xsrf-Protection"]; len(val) != 1 || val[0] != "1" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "CSRF protection error"})
			return
		}
	}

	ctx.Next()
}

func getCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-Xsrf-Protection, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
