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
	b.app.GET("/unauthorized", b.ViewUnauthorizedPage)
	b.app.GET("/error", b.ViewErrorPage)
	b.app.GET("/login", b.ViewLoginPage)

	// private routes
	b.app.GET("/", b.ViewHomePage)

	apiGroup := b.app.Group("/api")
	{
		notesGroup := apiGroup.Group("/notes")
		{
			notesGroup.GET("", b.GetAllNotes)
			notesGroup.GET(":id", b.GetNoteByID)
			notesGroup.POST("", b.CreateNote)
			notesGroup.PUT(":id", b.UpdateNote)
			notesGroup.DELETE(":id", b.DeleteNote)
		}
	}
}

// sendJSONOrHTML sends the provided status and either the JSON data or the HTML component
// based on the request's 'Accept' header. A value of 'application/json' will result in
// JSON, all others will result in a HTML response.
func sendJSONOrHTML(ctx *gin.Context, status int, data *gin.H, template templ.Component) {
	if ctx.GetHeader("Accept") == "application/json" {
		ctx.JSON(status, &data)
		return
	}

	template.Render(ctx, ctx.Writer)
}

// sendJSONOrRedirect sends the provided status and either the JSON data or a redirect
// to the provided target based on the request's 'Accept' header. A value of
// 'application/json' will result in JSON, all others will result in a redirect.
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
