package gin_binding

// import (
// 	"time"

// 	jwt "github.com/appleboy/gin-jwt/v2"
// 	"github.com/gin-gonic/gin"
// 	"github.com/oalexander6/passman/pkg/config"
// 	"github.com/oalexander6/passman/pkg/entities"
// )

// const (
// 	JWT_IDENTITY_KEY   = "userID"
// 	JWT_USER_EMAIL_KEY = "userEmail"
// )

// type claimsUserData struct {
// 	UserID entities.ID
// 	Email  string
// }

// func getJWTMiddleware(conf config.Config) *jwt.GinJWTMiddleware {
// 	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
// 		Realm:       "passman",
// 		Key:         []byte(conf.JWTSecretKey),
// 		Timeout:     time.Hour,
// 		MaxRefresh:  time.Hour,
// 		IdentityKey: JWT_IDENTITY_KEY,
// 		PayloadFunc: func(data interface{}) jwt.MapClaims {
// 			if v, ok := data.(*claimsUserData); ok {
// 				return jwt.MapClaims{
// 					JWT_IDENTITY_KEY:   v.UserID,
// 					JWT_USER_EMAIL_KEY: v.Email,
// 				}
// 			}
// 			return jwt.MapClaims{}
// 		},
// 		IdentityHandler: func(c *gin.Context) interface{} {
// 			claims := jwt.ExtractClaims(c)
// 			return &claimsUserData{
// 				UserID: claims[JWT_IDENTITY_KEY].(string),
// 			}
// 		},
// 		Authenticator: func(c *gin.Context) (interface{}, error) {
// 			var loginVals login
// 			if err := c.ShouldBind(&loginVals); err != nil {
// 				return "", jwt.ErrMissingLoginValues
// 			}
// 			userID := loginVals.Username
// 			password := loginVals.Password

// 			if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
// 				return &User{
// 					UserName:  userID,
// 					LastName:  "Bo-Yi",
// 					FirstName: "Wu",
// 				}, nil
// 			}

// 			return nil, jwt.ErrFailedAuthentication
// 		},
// 		Authorizator: func(data interface{}, c *gin.Context) bool {
// 			if v, ok := data.(*User); ok && v.UserName == "admin" {
// 				return true
// 			}

// 			return false
// 		},
// 		Unauthorized: func(c *gin.Context, code int, message string) {
// 			c.JSON(code, gin.H{
// 				"code":    code,
// 				"message": message,
// 			})
// 		},
// 		// TokenLookup is a string in the form of "<source>:<name>" that is used
// 		// to extract token from the request.
// 		// Optional. Default value "header:Authorization".
// 		// Possible values:
// 		// - "header:<name>"
// 		// - "query:<name>"
// 		// - "cookie:<name>"
// 		// - "param:<name>"
// 		TokenLookup: "header: Authorization, query: token, cookie: jwt",
// 		// TokenLookup: "query:token",
// 		// TokenLookup: "cookie:token",

// 		// TokenHeadName is a string in the header. Default value is "Bearer"
// 		TokenHeadName: "Bearer",

// 		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
// 		TimeFunc: time.Now,
// 	})
// }
