package httpserver

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oalexander6/passman/config"
	"github.com/oalexander6/passman/pkg/logger"
	"github.com/oalexander6/passman/pkg/models"
)

type Server struct {
	config *config.Config
	models *models.Models
	router http.Handler
}

// Initializes a new instance of a Gin HTTP server. If the environment set in the provided
// config is PROD, Gin will run in release mode, otherwise debug mode.
func New(conf *config.Config, store models.Store) *Server {
	if conf.Env == config.PROD_ENV {
		gin.SetMode(gin.ReleaseMode)
	}

	s := &Server{
		config: conf,
		models: models.New(store, conf),
	}

	s.router = s.createRouter()

	return s
}

// Runs the server. Gracefully shuts down on SIGINT or SIGTERM.
func (s *Server) Run() {
	srv := &http.Server{
		Addr:    ":" + s.config.Port,
		Handler: s.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal().Msgf("Error while serving: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info().Msg("Received shutdown signal")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal().Msgf("Server shutdown error: %s", err)
	}
	<-ctx.Done()
	logger.Log.Info().Msg("Server exiting")
}

// Adds a JSON response with a status code. Adds a timestamp and the request ID to the JSON body.
func json(ctx *gin.Context, status int, data gin.H) {
	data["timestamp"] = time.Now().Format(time.RFC3339)
	data["requestId"] = ctx.MustGet(requestIDKey)
	ctx.JSON(status, data)
}
