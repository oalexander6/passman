package gin_binding

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/oalexander6/passman/pkg/config"
	"github.com/oalexander6/passman/pkg/pages"
	"github.com/oalexander6/passman/pkg/services"
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

func (b *GinBinding) attachHandlers() {
	b.app.Static("assets", b.config.StaticFilePath)
	b.app.StaticFile("favicon.ico", fmt.Sprintf("%s/favicon.png", b.config.StaticFilePath))

	b.app.GET("/unauthorized", b.ViewUnauthorizedPage)
	b.app.GET("/error", b.ViewErrorPage)
	b.app.GET("/login", b.ViewLoginPage)

	authHandlers, _ := b.services.Auth.Handlers()
	auth := gin.WrapH(authHandlers)
	b.app.Any("/auth/*action", auth)

	authMiddleware := b.services.Auth.Middleware()
	requireAuth := gin.WrapH(authMiddleware.Auth(b.app.Handler()))
	b.app.Use(requireAuth)
	b.app.Use(func(ctx *gin.Context) {
		if ctx.Writer.Status() > 399 {
			ctx.Abort()
		}
	})

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
