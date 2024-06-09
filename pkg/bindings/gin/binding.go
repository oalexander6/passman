package gin_binding

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oalexander6/passman/pkg/services"
)

// GinBinding represents an instance of a Gin application and the associated configuration.
type GinBinding struct {
	services   *services.Services
	debugMode  bool
	secretKey  string
	listenAddr string
	app        *gin.Engine
}

func New(services *services.Services, listenAddr string, debugMode bool, secretKey string) *GinBinding {
	ginBinding := &GinBinding{
		services:   services,
		debugMode:  debugMode,
		secretKey:  secretKey,
		listenAddr: listenAddr,
	}

	if debugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.New()
	app.SetTrustedProxies(nil)
	app.Use(gin.Recovery())
	app.Use(gin.Logger())

	ginBinding.app = app

	ginBinding.attachHandlers()

	return ginBinding
}

// Run starts the application with the provided configuration.
// It returns an error if the application crashes.
func (b *GinBinding) Run() error {
	return b.app.Run(b.listenAddr)
}

func (b *GinBinding) attachHandlers() {
	b.app.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})

	notesGroup := b.app.Group("/notes")
	{
		notesGroup.GET("", b.GetAllNotes)
		notesGroup.GET(":id", b.GetNoteByID)
		notesGroup.POST("", b.CreateNote)
		notesGroup.PUT(":id", b.UpdateNote)
		notesGroup.DELETE(":id", b.DeleteNote)
	}
}
