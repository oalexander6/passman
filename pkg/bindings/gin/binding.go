package gin_binding

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oalexander6/passman/pkg/entities"
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

// attachHandlers adds the handlers to the underlying Gin app.
func (b *GinBinding) attachHandlers() {
	b.app.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})

	b.app.GET("/notes", func(ctx *gin.Context) {
		notes, err := b.services.GetAllNotes(ctx)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"notes": notes,
		})
	})

	b.app.GET("/notes/:id", func(ctx *gin.Context) {
		noteID := ctx.Param("id")

		note, err := b.services.GetNoteByID(ctx, noteID)
		if err != nil {
			if errors.Is(err, entities.ErrNotFound) {
				ctx.AbortWithStatus(http.StatusNotFound)
				return
			}
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"note": note,
		})
	})

	b.app.POST("/notes", func(ctx *gin.Context) {
		var noteInput entities.NoteInput

		if err := ctx.ShouldBind(&noteInput); err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		savedNote, err := b.services.CreateNote(ctx, noteInput)
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"note": savedNote,
		})
	})
}
