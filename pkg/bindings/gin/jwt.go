package gin_binding

import (
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/oalexander6/passman/pkg/entities"
)

const (
	JWT_IDENTITY_KEY = "accountID"
)

type claimsAccountData struct {
	AccountID entities.ID
}

func (b GinBinding) initJWTParams() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:       "passman",
		Key:         []byte(b.config.JWTSecretKey),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: JWT_IDENTITY_KEY,
		SendCookie:  true,

		PayloadFunc:     getAccountClaimsFromAccount,
		IdentityHandler: getAccountIDFromClaims,
		Authenticator:   b.handleAuthentication,
		Authorizator:    b.handleAuthorization,
		Unauthorized:    handleUnauthorizedError,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	}
}

func (b *GinBinding) getJWTProvider() *jwt.GinJWTMiddleware {
	jwtMiddleware, err := jwt.New(b.initJWTParams())
	if err != nil {
		panic(err)
	}

	return jwtMiddleware
}

func (b *GinBinding) handleAuthentication(c *gin.Context) (interface{}, error) {
	var input entities.AccountLoginInput
	if err := c.ShouldBind(&input); err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	account, err := b.services.Login(c, input)
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	return &claimsAccountData{
		AccountID: account.ID,
	}, nil
}

func (b *GinBinding) handleAuthorization(data interface{}, c *gin.Context) bool {
	return true
}

func handleUnauthorizedError(c *gin.Context, code int, message string) {
	sendJSONOrRedirect(c, http.StatusUnauthorized, &gin.H{"message": "not authorized"}, "/login")
}

func getAccountIDFromClaims(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &claimsAccountData{
		AccountID: claims[JWT_IDENTITY_KEY].(string),
	}
}

func getAccountClaimsFromAccount(data interface{}) jwt.MapClaims {
	if v, ok := data.(*entities.Account); ok {
		return jwt.MapClaims{
			JWT_IDENTITY_KEY: v.ID,
		}
	}
	return jwt.MapClaims{}
}

func extractAccountID(c *gin.Context) *claimsAccountData {
	claims := jwt.ExtractClaims(c)

	return &claimsAccountData{
		AccountID: claims[JWT_IDENTITY_KEY].(entities.ID),
	}
}
