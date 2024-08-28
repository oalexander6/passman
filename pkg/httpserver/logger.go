package httpserver

import (
	"net/http"
	"time"

	"github.com/oalexander6/passman/pkg/logger"
	"github.com/urfave/negroni"
)

func logMiddleware(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	lrw := negroni.NewResponseWriter(rw)
	startTime := time.Now()

	next(lrw, r)

	statusCode := lrw.Status()
	size := lrw.Size()
	logger.Log.Info().Msgf("%s %s | %d %s | %db %dms", r.Method, r.URL.Path, statusCode, http.StatusText(statusCode), size, time.Since(startTime).Milliseconds())
}
