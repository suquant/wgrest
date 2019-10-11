// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/swag"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/suquant/wgrest/restapi/operations"
	"github.com/suquant/wgrest/restapi/operations/wireguard"
	wireguardHandlers "github.com/suquant/wgrest/wireguard"
)

//go:generate swagger generate server --target ../../wgrest --name Wgrest --spec ../swagger.yml

var authFlags = struct {
	Token string `long:"token" description:"authentication token"`
}{}

func configureFlags(api *operations.WgrestAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		swag.CommandLineOptionsGroup{
			ShortDescription: "Auth flags",
			LongDescription:  "Authentication flags",
			Options:          &authFlags,
		},
	}
}

func configureAPI(api *operations.WgrestAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "Token" header is set
	api.KeyAuth = func(token string) (interface{}, error) {
		authToken := authFlags.Token
		if authToken == "" {
			// pass, no auth required
			return "", nil
		}

		if authToken == token {
			return "", nil
		}

		return nil, nil
	}

	wireguardHandlers.ApplyAPI(api)

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()
	if api.WireguardDeviceCreateHandler == nil {
		api.WireguardDeviceCreateHandler = wireguard.DeviceCreateHandlerFunc(func(params wireguard.DeviceCreateParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation wireguard.DeviceCreate has not yet been implemented")
		})
	}
	if api.WireguardDeviceDeleteHandler == nil {
		api.WireguardDeviceDeleteHandler = wireguard.DeviceDeleteHandlerFunc(func(params wireguard.DeviceDeleteParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation wireguard.DeviceDelete has not yet been implemented")
		})
	}
	if api.WireguardDeviceGetHandler == nil {
		api.WireguardDeviceGetHandler = wireguard.DeviceGetHandlerFunc(func(params wireguard.DeviceGetParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation wireguard.DeviceGet has not yet been implemented")
		})
	}
	if api.WireguardDeviceListHandler == nil {
		api.WireguardDeviceListHandler = wireguard.DeviceListHandlerFunc(func(params wireguard.DeviceListParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation wireguard.DeviceList has not yet been implemented")
		})
	}
	if api.WireguardPeerCreateHandler == nil {
		api.WireguardPeerCreateHandler = wireguard.PeerCreateHandlerFunc(func(params wireguard.PeerCreateParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation wireguard.PeerCreate has not yet been implemented")
		})
	}
	if api.WireguardPeerDeleteHandler == nil {
		api.WireguardPeerDeleteHandler = wireguard.PeerDeleteHandlerFunc(func(params wireguard.PeerDeleteParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation wireguard.PeerDelete has not yet been implemented")
		})
	}
	if api.WireguardPeerGetHandler == nil {
		api.WireguardPeerGetHandler = wireguard.PeerGetHandlerFunc(func(params wireguard.PeerGetParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation wireguard.PeerGet has not yet been implemented")
		})
	}
	if api.WireguardPeerListHandler == nil {
		api.WireguardPeerListHandler = wireguard.PeerListHandlerFunc(func(params wireguard.PeerListParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation wireguard.PeerList has not yet been implemented")
		})
	}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
