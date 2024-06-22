package gin_binding

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
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

	ginBinding.setupAuthMiddleware()
	ginBinding.setupCSRFMiddleware()
	ginBinding.attachHandlers()

	return ginBinding
}

// Run starts the application with the provided configuration.
// It returns an error if the application crashes.
func (b *GinBinding) Run() error {
	return b.app.Run(b.config.Host)
}

func (b *GinBinding) attachHandlers() {
	b.app.Static("assets", b.config.StaticFilePath)
	b.app.StaticFile("favicon.ico", fmt.Sprintf("%s/favicon.png", b.config.StaticFilePath))

	b.app.GET("/unauthorized", b.ViewUnauthorizedPage)
	b.app.GET("/error", b.ViewErrorPage)
	b.app.GET("/", b.ViewHomePage)

	notesGroup := b.app.Group("/notes")
	{
		notesGroup.GET("", b.GetAllNotes)
		notesGroup.GET(":id", b.GetNoteByID)
		notesGroup.POST("", b.CreateNote)
		notesGroup.PUT(":id", b.UpdateNote)
		notesGroup.DELETE(":id", b.DeleteNote)
	}
}

func (b *GinBinding) setupAuthMiddleware() {
	store := cookie.NewStore([]byte(b.config.SecretKey))
	store.Options(sessions.Options{
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	b.app.Use(sessions.Sessions(b.config.CookieName, store))
}

func (b *GinBinding) setupCSRFMiddleware() {
	csrfIgnoreMethods := []string{"GET", "HEAD", "OPTIONS"}

	if !b.config.UseCSRFTokens {
		csrfIgnoreMethods = append(csrfIgnoreMethods, "POST")
	}

	b.app.Use(csrf.Middleware(csrf.Options{
		IgnoreMethods: csrfIgnoreMethods,
		Secret:        b.config.CSRFSecret,
		ErrorFunc: func(c *gin.Context) {
			sendJSONOrRedirect(
				c,
				http.StatusBadRequest,
				&gin.H{},
				"/unauthorized",
			)
		},
	}))
}

func sendJSONOrHTML(ctx *gin.Context, status int, data *gin.H, template templ.Component) {
	if ctx.GetHeader("Accept") == "application/json" {
		ctx.JSON(status, &data)
		return
	}

	template.Render(ctx, ctx.Writer)
}

func sendJSONOrRedirect(ctx *gin.Context, status int, data *gin.H, target string) {
	if ctx.GetHeader("Accept") == "application/json" {
		ctx.JSON(status, &data)
		return
	}

	ctx.Redirect(http.StatusFound, target)
}

func (b *GinBinding) ViewUnauthorizedPage(ctx *gin.Context) {
	sendJSONOrHTML(ctx, http.StatusUnauthorized, &gin.H{}, pages.NotAuthorized())
}

func (b *GinBinding) ViewErrorPage(ctx *gin.Context) {
	sendJSONOrHTML(ctx, http.StatusInternalServerError, &gin.H{}, pages.Error())
}
