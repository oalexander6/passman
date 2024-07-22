package gin_binding

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/oalexander6/passman/pkg/entities"
)

const (
	SESSION_USER_ID_KEY    = "userId"
	SESSION_USER_EMAIL_KEY = "userEmail"
	SESSION_LOGGED_IN_KEY  = "isLoggedIn"
)

type sessionUserData struct {
	UserID     entities.ID
	Email      string
	IsLoggedIn bool
}

func requireAuthMiddleware(ctx *gin.Context) {
	sessionUserData := getSessionUser(ctx)

	if !sessionUserData.IsLoggedIn {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
}

func setSessionUser(ctx *gin.Context, id string, email string) {
	session := sessions.Default(ctx)
	session.Set(SESSION_USER_ID_KEY, id)
	session.Set(SESSION_USER_EMAIL_KEY, email)
	session.Set(SESSION_LOGGED_IN_KEY, true)
	session.Save()
}

func getSessionUser(ctx *gin.Context) sessionUserData {
	session := sessions.Default(ctx)

	userID, ok := session.Get(SESSION_LOGGED_IN_KEY).(entities.ID)
	if !ok {
		return sessionUserData{}
	}

	userEmail, ok := session.Get(SESSION_USER_EMAIL_KEY).(string)
	if !ok {
		return sessionUserData{}
	}

	isLoggedIn, ok := session.Get(SESSION_LOGGED_IN_KEY).(bool)
	if !ok {
		return sessionUserData{}
	}

	return sessionUserData{
		IsLoggedIn: isLoggedIn,
		UserID:     userID,
		Email:      userEmail,
	}
}

func removeSessionUser(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
}
