package gin_binding

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oalexander6/passman/pkg/pages"
)

func (b *GinBinding) ViewLoginPage(ctx *gin.Context) {
	sendJSONOrHTML(ctx, http.StatusOK, &gin.H{}, pages.Login())
}
