// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	"github.com/hallanneves/autenticador/autenticador"
	"github.com/hallanneves/autenticador/conf"
	"github.com/hallanneves/autenticador/models"
	"github.com/hallanneves/autenticador/restapi/operations"
	"github.com/hallanneves/autenticador/restapi/operations/auth"
)

//go:generate swagger generate server --target ../../autenticador --name Autenticador --spec ../swagger.yaml
var autenticadorFlags = struct {
	ConfigFile string `long:"ConfigFile" description:"Arquivo de configuracao padrao (default: conf/conf.json)" default:"conf/conf.json"`
}{}

func configureFlags(api *operations.AutenticadorAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{
		swag.CommandLineOptionsGroup{
			ShortDescription: "Autenticador Flags",
			LongDescription:  "",
			Options:          &autenticadorFlags,
		},
	}
}

func configureAPI(api *operations.AutenticadorAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.AuthValidaCredenciaisHandler = auth.ValidaCredenciaisHandlerFunc(func(params auth.ValidaCredenciaisParams) middleware.Responder {
		status, err := autenticador.ValidaAutenticacao(params.Body)
		if err == nil {
			if status == 200 {
				return auth.NewValidaCredenciaisOK()
			} else if status == 401 {
				return auth.NewValidaCredenciaisUnauthorized()
			}
		}
		mensagem := models.Erro{Mensagem: err.Error()}
		log.Println(err.Error())
		return auth.NewValidaCredenciaisInternalServerError().WithPayload(&mensagem)
	})

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

	//* Le o arquivo de configuracao
	var err error
	err = conf.LerConfig(autenticadorFlags.ConfigFile)
	if err != nil {
		log.Println("Erro leitura de arquivo de configuracao: " + err.Error())
		os.Exit(255)
	}

	//* Inicializa o Mysql
	err = autenticador.InicializaMysql()
	if err != nil {
		log.Println("Erro conexao com Mysql: " + err.Error())
		os.Exit(255)
	}

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
