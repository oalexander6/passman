package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oalexander6/passman/config"
	"github.com/oalexander6/passman/pkg/models"
)

type Server struct {
	config *config.Config
	models *models.Models
	server *gin.Engine
}

// Initializes a new instance of a Gin HTTP server. If the environment set in the provided
// config is PROD, Gin will run in release mode, otherwise debug mode.
func New(conf *config.Config, store models.Store) *Server {
	if conf.Env == config.PROD_ENV {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.SetTrustedProxies(nil)

	s := &Server{
		config: conf,
		models: models.New(store, conf),
		server: r,
	}

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello")
	})
	r.GET("/notes", s.handleGetAllNotes)

	return s
}

func (s *Server) Run() error {
	return s.server.Run(":" + s.config.Port)
}

func (s *Server) handleGetAllNotes(ctx *gin.Context) {
	notes, err := s.models.NoteGetAll(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"notes": notes})
}
