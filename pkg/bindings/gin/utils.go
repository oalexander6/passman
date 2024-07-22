package gin_binding

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

func (b *GinBinding) getCSRFToken(c *gin.Context) string {
	csrfToken := ""
	if b.config.UseCSRFTokens {
		csrfToken = csrf.GetToken(c)
	}
	return csrfToken
}

// sendJSONOrHTML sends the provided status and either the JSON data or the HTML component
// based on the request's 'Accept' header. A value of 'application/json' will result in
// JSON, all others will result in a HTML response.
func sendJSONOrHTML(c *gin.Context, status int, data *gin.H, template templ.Component) {
	if c.GetHeader("Accept") == "application/json" {
		c.JSON(status, &data)
		return
	}

	template.Render(c, c.Writer)
}

// sendJSONOrRedirect sends the provided status and either the JSON data or a redirect
// to the provided target based on the request's 'Accept' header. A value of
// 'application/json' will result in JSON, all others will result in a redirect.
func sendJSONOrRedirect(c *gin.Context, status int, data *gin.H, target string) {
	if c.GetHeader("Accept") == "application/json" {
		c.JSON(status, &data)
		return
	}

	c.Redirect(http.StatusFound, target)
}
