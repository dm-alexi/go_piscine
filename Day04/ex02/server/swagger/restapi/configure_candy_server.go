// This file is safe to edit. Once it exists it will not be overwritten

package restapi

// #cgo CFLAGS: -g -Wall
//#include <stdlib.h>
//#include "cow.h"
import "C"
import (
	"crypto/tls"
	"fmt"
	"net/http"
	"unsafe"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/dm-alexi/go_piscine/Day04/ex02/server/swagger/restapi/operations"
)

//go:generate swagger generate server --target ../../swagger --name CandyServer --spec ../config.yml --principal interface{} --exclude-main

func configureFlags(api *operations.CandyServerAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.CandyServerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	pricelist := map[string]int64{"CE": 10, "AA": 15, "NT": 17, "DE": 21, "YR": 23}
	api.BuyCandyHandler = operations.BuyCandyHandlerFunc(
		func(params operations.BuyCandyParams) middleware.Responder {
			ptr := C.ask_cow(C.CString("Thank you!"))
			defer C.free(unsafe.Pointer(ptr))
			cowSay := C.GoString(ptr)
			if pricelist[*params.Order.CandyType] == 0 || *params.Order.CandyCount <= 0 {
				response := operations.BuyCandyBadRequestBody{Error: "Invalid candy type"}
				if *params.Order.CandyCount <= 0 {
					response.Error = "Invalid candy count"
				}
				return operations.NewBuyCandyBadRequest().WithPayload(&response)
			}
			change := *params.Order.Money - *params.Order.CandyCount*pricelist[*params.Order.CandyType]
			if change < 0 {
				response := operations.BuyCandyPaymentRequiredBody{Error: fmt.Sprintf("You need %d more money!", -change)}
				return operations.NewBuyCandyPaymentRequired().WithPayload(&response)
			}
			response := operations.BuyCandyCreatedBody{Thanks: cowSay, Change: &change}
			return operations.NewBuyCandyCreated().WithPayload(&response)
		})

	if api.BuyCandyHandler == nil {
		api.BuyCandyHandler = operations.BuyCandyHandlerFunc(func(params operations.BuyCandyParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.BuyCandy has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
	tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
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
