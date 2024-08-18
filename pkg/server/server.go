package server

import (
	"net/http"

	"github.com/oalexander6/passman/config"
	"github.com/oalexander6/passman/pkg/logger"
	"github.com/oalexander6/passman/pkg/models"
)

type Server struct {
	config *config.Config
	models *models.Models
	server *http.ServeMux
}

func New(conf *config.Config, store models.Store) *Server {
	mux := http.NewServeMux()

	return &Server{
		config: conf,
		models: models.New(store, conf),
		server: mux,
	}
}

func (s *Server) Run() error {
	httpServer := &http.Server{
		Addr:    ":" + s.config.Port,
		Handler: s.server,
	}

	logger.Log.Info().Msgf("listening on %s\n", httpServer.Addr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Info().Msgf("error listening and serving: %s\n", err)
	}

	return nil
}
