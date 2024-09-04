package httpserver

import (
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/oalexander6/passman/config"
	"github.com/oalexander6/passman/pkg/logger"
)

func (s *Server) createRouter() *gin.Engine {
	r := gin.New()
	r.SetTrustedProxies(nil)

	r.Use(gin.Recovery())
	r.Use(gin.LoggerWithWriter(logger.Log))

	if s.config.Env != config.LOCAL_ENV {
		r.Use(static.Serve("/", static.LocalFile("./web/dist", true)))
	}

	r.Use(requestIDMiddleware)
	r.Use(getSecurityHeadersMiddleware(s.config.AllowedHost))
	r.Use(csrfHeaderMiddleware)

	if s.config.Env == config.LOCAL_ENV {
		r.Use(getCORSMiddleware())
	}

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/", s.hello)
		apiGroup.GET("/notes", s.handleGetAllNotes)
		apiGroup.POST("/notes", s.handleCreateNote)
	}

	return r
}

func (s *Server) hello(ctx *gin.Context) {
	ctx.String(http.StatusOK, "Ok")
}
