package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oalexander6/passman/config"
	"github.com/oalexander6/passman/pkg/logger"
)

func (s *Server) createRouter() *gin.Engine {
	r := gin.New()
	r.SetTrustedProxies(nil)

	r.Use(gin.Recovery())
	r.Use(gin.LoggerWithWriter(logger.Log))
	r.Use(requestIDMiddleware)
	r.Use(getSecurityHeadersMiddleware(s.config.AllowedHost))
	r.Use(csrfHeaderMiddleware)

	if s.config.Env == config.LOCAL_ENV {
		r.Use(getCORSMiddleware())
	}

	r.GET("/", s.hello)
	r.GET("/notes", s.handleGetAllNotes)
	r.POST("/notes", s.handleCreateNote)

	return r
}

func (s *Server) hello(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Ok")
}
