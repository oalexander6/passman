package httpserver

import (
	"net/http"

	"github.com/oalexander6/passman/config"
	"github.com/oalexander6/passman/pkg/logger"
	"github.com/oalexander6/passman/pkg/models"
	"github.com/urfave/negroni"
)

type Server struct {
	config *config.Config
	models *models.Models
	server http.Handler
}

func New(conf *config.Config, store models.Store) *Server {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	mw := negroni.New()
	mw.Use(negroni.NewRecovery())
	mw.Use(negroni.HandlerFunc(logMiddleware))
	mw.UseHandler(mux)

	return &Server{
		config: conf,
		models: models.New(store, conf),
		server: mw,
	}
}

func (s *Server) Run() error {
	logger.Log.Info().Msgf("listening on %s\n", s.config.Port)
	if err := http.ListenAndServe(":"+s.config.Port, s.server); err != nil && err != http.ErrServerClosed {
		logger.Log.Info().Msgf("error listening and serving: %s\n", err)
	}

	return nil
}
