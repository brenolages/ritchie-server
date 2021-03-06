package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"ritchie-server/server"
	"ritchie-server/server/starter"
)

var h Handler

type Handler struct {
	LoginHandler            server.DefaultHandler
	KeycloakHandler         server.DefaultHandler
	UserHandler             server.DefaultHandler
	CredentialConfigHandler server.DefaultHandler
	ConfigHealth            server.DefaultHandler
	OauthHandler            server.DefaultHandler
	UsageLoggerHandler      server.DefaultHandler
	CliVersionHandler       server.DefaultHandler
	RepositoryHandler       server.DefaultHandler
	MiddlewareHandler       server.MiddlewareHandler
	CredentialHandler       server.CredentialHandler
}

func init() {
	i, err := starter.NewConfiguration()
	if err != nil {
		log.Fatalf("Failed to load server configuration: %v", err)
	}
	h = Handler{
		LoginHandler:            i.LoadLoginHandler(),
		KeycloakHandler:         i.LoadConfigHandler(),
		UserHandler:             i.LoadUserHandler(),
		CredentialConfigHandler: i.LoadCredentialConfigHandler(),
		ConfigHealth:            i.LoadConfigHealth(),
		OauthHandler:            i.LoadOauthHandler(),
		UsageLoggerHandler:      i.LoadUsageLoggerHandler(),
		CliVersionHandler:       i.LoadCliVersionHandler(),
		RepositoryHandler:       i.LoadRepositoryHandler(),
		MiddlewareHandler:       i.LoadMiddlewareHandler(),
		CredentialHandler:       i.LoadCredentialHandler(),
	}
}

func main() {

	log.Info("Starting server")
	http.Handle("/login", h.MiddlewareHandler.Filter(h.LoginHandler.Handler()))
	http.Handle("/keycloak", h.MiddlewareHandler.Filter(h.KeycloakHandler.Handler()))
	http.Handle("/users", h.MiddlewareHandler.Filter(h.UserHandler.Handler()))
	http.Handle("/credentials/admin", h.MiddlewareHandler.Filter(h.CredentialHandler.HandlerAdmin()))
	http.Handle("/credentials/me", h.MiddlewareHandler.Filter(h.CredentialHandler.HandlerMe()))
	http.Handle("/credentials/me/", h.MiddlewareHandler.Filter(h.CredentialHandler.HandlerMe()))
	http.Handle("/credentials/config", h.MiddlewareHandler.Filter(h.CredentialConfigHandler.Handler()))
	http.Handle("/metrics", h.MiddlewareHandler.Filter(promhttp.Handler()))
	http.Handle("/usage", h.MiddlewareHandler.Filter(h.UsageLoggerHandler.Handler()))
	http.Handle("/health", h.MiddlewareHandler.Filter(h.ConfigHealth.Handler()))
	http.Handle("/oauth", h.MiddlewareHandler.Filter(h.OauthHandler.Handler()))
	http.Handle("/cli-version", h.MiddlewareHandler.Filter(h.CliVersionHandler.Handler()))
	http.Handle("/repositories", h.MiddlewareHandler.Filter(h.RepositoryHandler.Handler()))
	log.Fatal(http.ListenAndServe(":3000", nil))
}
