package httpserver

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Adds a JSON response with a status code. Adds a timestamp and the request ID to the JSON body.
func json(ctx *gin.Context, status int, data gin.H) {
	data["timestamp"] = time.Now().Format(time.RFC3339)
	data["requestId"] = ctx.MustGet(requestIDKey)
	ctx.JSON(status, data)
}
