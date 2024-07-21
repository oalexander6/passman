package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	"github.com/go-pkgz/auth/logger"
	"github.com/go-pkgz/auth/provider"
	"github.com/go-pkgz/auth/token"
	"github.com/oalexander6/passman/pkg/config"
)

// newAuthService intializes an instance of a configured auth provider.
func newAuthService(conf *config.Config) *auth.Service {
	protocol := "http"
	if conf.UseSSL {
		protocol = "https"
	}

	options := auth.Opts{
		SecretReader: token.SecretFunc(func(id string) (string, error) {
			return conf.SecretKey, nil
		}),

		DisableXSRF: !conf.UseCSRFTokens,

		TokenDuration:  time.Minute * 30,
		CookieDuration: time.Hour * 24,
		SameSiteCookie: http.SameSiteStrictMode,
		SecureCookies:  true,
		Issuer:         "passman",

		URL:         fmt.Sprintf("%s://%s", protocol, conf.Host),
		AvatarStore: avatar.NewNoOp(),

		Logger: logger.Std,
	}

	service := auth.NewService(options)
	service.AddDirectProvider("local", provider.CredCheckerFunc(func(user string, password string) (bool, error) {
		return true, nil
	}))

	// srv := initGoauth2Srv()
	// sopts := provider.CustomServerOpt{
	// 	URL:           "http://127.0.0.1:9096",
	// 	L:             options.Logger,
	// 	WithLoginPage: true,
	// }
	// // create custom provider and prepare params for handler
	// prov := provider.NewCustomServer(srv, sopts)

	// // Start server
	// go prov.Run(context.Background())
	// service.AddCustomProvider("custom123", auth.Client{Cid: "cid", Csecret: "csecret"}, prov.HandlerOpt)

	return service
}

// func initGoauth2Srv() *goauth2.Server {
// 	manager := manage.NewDefaultManager()
// 	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

// 	// token store
// 	manager.MustTokenStorage(store.NewMemoryTokenStore())

// 	// generate jwt access token
// 	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("custom", []byte("00000000"), jwt.SigningMethodHS512))

// 	// client memory store
// 	clientStore := store.NewClientStore()
// 	err := clientStore.Set("cid", &models.Client{
// 		ID:     "cid",
// 		Secret: "csecret",
// 		Domain: "http://127.0.0.1:8081",
// 	})
// 	if err != nil {
// 		log.Printf("failed to set up a client store for go-oauth2/oauth2 server, %s", err)
// 	}
// 	manager.MapClientStorage(clientStore)

// 	srv := goauth2.NewServer(goauth2.NewConfig(), manager)

// 	srv.SetUserAuthorizationHandler(func(_ http.ResponseWriter, r *http.Request) (string, error) {
// 		if r.Form.Get("username") != "admin" || r.Form.Get("password") != "admin" {
// 			return "", fmt.Errorf("wrong creds. Use: admin admin")
// 		}
// 		return "custom123_admin", nil
// 	})

// 	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
// 		log.Printf("Internal Error: %s", err.Error())
// 		return
// 	})

// 	srv.SetResponseErrorHandler(func(re *errors.Response) {
// 		log.Printf("Response Error: %s", re.Error.Error())
// 	})

// 	return srv
// }
