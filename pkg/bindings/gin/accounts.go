package gin_binding

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oalexander6/passman/pkg/entities"
	"github.com/oalexander6/passman/pkg/pages"
)

func (b *GinBinding) ViewLoginPage(c *gin.Context) {
	pages.Login(b.getCSRFToken(c)).Render(c, c.Writer)
}

func (b *GinBinding) HandleLogin(c *gin.Context) {
	var input entities.AccountLoginInput
	if err := c.ShouldBind(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	account, err := b.services.Login(c, input)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	setSessionUser(c, account.ID, account.Email)

	c.Redirect(http.StatusOK, "/")
}

func (b *GinBinding) HandleRegister(c *gin.Context) {
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

	setSessionUser(c, account.ID, account.Email)

	sendJSONOrRedirect(c, http.StatusOK, &gin.H{"account": account}, "/")
}

func (b *GinBinding) HandleLogout(c *gin.Context) {
	removeSessionUser(c)

	sendJSONOrRedirect(c, http.StatusOK, &gin.H{"message": "logged out"}, "/")
}

func (b *GinBinding) HandleStatus(c *gin.Context) {
	sessionUser := getSessionUser(c)

	c.JSON(http.StatusOK, gin.H{"user": sessionUser})
}
