package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) createRouter() *gin.Engine {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.Use(requestIDMiddleware)
	r.Use(getSecurityHeadersMiddleware(s.config.AllowedHost))

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Ok")
	})
	r.GET("/notes", s.handleGetAllNotes)
	r.POST("/notes", s.handleCreateNote)

	return r
}
