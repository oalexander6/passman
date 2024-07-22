package gin_binding

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/oalexander6/passman/pkg/entities"
	"github.com/oalexander6/passman/pkg/pages"
)

func (b *GinBinding) viewLoginPage(c *gin.Context) {
	csrfToken := b.getCSRFToken(c)
	sendJSONOrHTML(c, http.StatusOK, &gin.H{"csrfToken": csrfToken}, pages.Login(csrfToken))
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

	session := sessions.Default(c)
	session.Set("user", account.ID)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	sendJSONOrRedirect(c, http.StatusOK, &gin.H{"account": account}, "/login")
}

func (b *GinBinding) handleStatus(c *gin.Context) {
	// id := extractAccountID(c)
	session := sessions.Default(c)
	user := session.Get("user")

	c.JSON(http.StatusOK, gin.H{"sessionUserID": user})
}
