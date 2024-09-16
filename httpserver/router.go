package httpserver

import (
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/oalexander6/passman/config"
	"github.com/oalexander6/passman/logger"
	"github.com/oalexander6/passman/models"
)

func (s *Server) createRouter(m models.Models) *gin.Engine {
	r := gin.New()
	r.SetTrustedProxies(nil)

	r.Use(gin.Recovery())
	r.Use(gin.LoggerWithWriter(logger.Log))

	if s.config.Env != config.LOCAL_ENV {
		r.Use(static.Serve("/", static.LocalFile("./web/dist", true)))
	}

	r.Use(requestIDMiddleware)
	r.Use(getSecurityHeadersMiddleware())
	r.Use(csrfHeaderMiddleware)

	apiGroup := r.Group("/api/v1")
	{
		apiGroup.GET("", HandleHello())
		apiGroup.GET("/notes", HandleGetAllNotes(m))
		apiGroup.POST("/notes", HandleCreateNote(m))
	}

	return r
}
