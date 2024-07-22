package gin_binding

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oalexander6/passman/pkg/entities"
	"github.com/oalexander6/passman/pkg/pages"
)

func (b *GinBinding) viewLoginPage(c *gin.Context) {
	pages.Login(b.getCSRFToken(c)).Render(c, c.Writer)
}

func (b *GinBinding) handleRegister(c *gin.Context) {
	var input entities.AccountInput
	if err := c.ShouldBind(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	account, err := b.services.Register(c, input)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	sendJSONOrRedirect(c, http.StatusOK, &gin.H{"account": account}, "/login")
}

func (b *GinBinding) handleStatus(c *gin.Context) {
	id := extractAccountID(c)

	c.JSON(http.StatusOK, gin.H{"id": id})
}
