package gin_binding

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/oalexander6/passman/pkg/config"
	"github.com/oalexander6/passman/pkg/pages"
	"github.com/oalexander6/passman/pkg/services"
	csrf "github.com/utrack/gin-csrf"
)

// GinBinding represents an instance of a Gin application and the associated configuration.
type GinBinding struct {
	services *services.Services
	config   *config.Config
	app      *gin.Engine
}

func New(services *services.Services, conf *config.Config) *GinBinding {
	ginBinding := &GinBinding{
		services: services,
		config:   conf,
	}

	if conf.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.New()
	app.SetTrustedProxies(nil)
	app.Use(gin.Recovery())
	app.Use(gin.Logger())

	ginBinding.app = app

	ginBinding.attachMiddleware()
	ginBinding.attachHandlers()

	return ginBinding
}

// Run starts the application with the provided configuration.
// It returns an error if the application crashes.
func (b *GinBinding) Run() error {
	return b.app.Run(b.config.Host)
}

func (b *GinBinding) attachMiddleware() {
	store := cookie.NewStore([]byte(b.config.JWTSecretKey))
	store.Options(sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Domain:   b.config.Host,
	})

	b.app.Use(sessions.Sessions(b.config.JWTCookieName, store))

	if b.config.UseCSRFTokens {
		b.app.Use(csrf.Middleware(csrf.Options{
			IgnoreMethods: []string{"GET", "HEAD", "OPTIONS"},
			Secret:        b.config.CSRFSecret,
			ErrorFunc: func(c *gin.Context) {
				c.AbortWithStatus(http.StatusBadRequest)
			},
		}))
	}
}

func (b *GinBinding) attachHandlers() {
	b.app.Static("assets", b.config.StaticFilePath)
	b.app.StaticFile("favicon.ico", fmt.Sprintf("%s/favicon.png", b.config.StaticFilePath))

	// public routes
	b.app.GET("/error", b.viewErrorPage)
	b.app.GET("/login", b.ViewLoginPage)
	b.app.GET("/register", b.ViewLoginPage)
	b.app.POST("/api/auth/login", b.HandleLogin)
	b.app.POST("/api/auth/register", b.HandleRegister)
	b.app.GET("/api/auth/logout", b.HandleLogout)
	b.app.GET("/api/auth/status", b.HandleStatus)

	// private routes
	b.app.Use(requireAuthMiddleware)
	b.app.GET("/", b.viewHomePage)

	apiGroup := b.app.Group("/api")
	{
		notesGroup := apiGroup.Group("/notes")
		{
			// notesGroup.GET("", b.HandleGetAllNotes)
			notesGroup.GET(":id", b.HandleGetNoteByID)
			notesGroup.POST("", b.HandleCreateNote)
			notesGroup.PUT(":id", b.HandleUpdateNote)
			notesGroup.DELETE(":id", b.HandleDeleteNote)
		}
	}
}

func (b *GinBinding) viewErrorPage(c *gin.Context) {
	sendJSONOrHTML(c, http.StatusInternalServerError, &gin.H{}, pages.Error())
}

func (b *GinBinding) viewHomePage(c *gin.Context) {
	notes, err := b.services.GetAllNotes(c)
	if err != nil {
		sendJSONOrHTML(c, http.StatusInternalServerError, &gin.H{}, pages.Error())
		return
	}

	sendJSONOrHTML(c, http.StatusOK, &gin.H{"notes": notes}, pages.Dashboard(notes, b.getCSRFToken(c)))
}
