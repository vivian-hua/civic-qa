package main

import (
	// standard
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	// 3rd party
	"github.com/gorilla/mux"

	// common
	"github.com/team-ravl/civic-qa/services/common/config"
	commonMiddleware "github.com/team-ravl/civic-qa/services/common/middleware"

	//pkg
	aggregator "github.com/team-ravl/civic-qa/services/logAggregator/pkg/middleware"

	// internal
	"github.com/team-ravl/civic-qa/services/gateway/internal/middleware"
	"github.com/team-ravl/civic-qa/services/gateway/internal/proxy"
)

const (
	// VersionBase is the API route base
	VersionBase = "/v0"
	// APIVersion is the API semantic version
	APIVersion = "v0.0.0"
)

var (
	// LoggingOutput is a file that recieves log outputs
	LoggingOutput = os.Stdout
)

func main() {
	// config
	var cfg config.Provider = &config.EnvProvider{}
	cfg.SetVerbose(true)

	// subservices
	accountService := cfg.GetOrFallback("ACCOUNT_SVC", "http://localhost:8080")
	formService := cfg.GetOrFallback("FORM_SVC", "http://localhost:7070")
	mailtoService := cfg.GetOrFallback("MAILTO_SVC", "http://localhost:9090")

	// Routers
	// base router
	router := mux.NewRouter()
	// API routers
	api := router.PathPrefix(VersionBase).Subrouter()
	// API routes requiring auth
	apiAuth := router.PathPrefix(VersionBase).Subrouter()

	// Middleware
	router.Use(middleware.NewCorrelatorMiddleware)
	router.Use(middleware.CorsMiddleware)
	router.Use(commonMiddleware.NewLoggingMiddleware(LoggingOutput))
	router.Use(aggregator.NewAggregatorMiddleware(&aggregator.Config{
		AggregatorAddress: cfg.GetOrFallback("AGG_ADDR", "http://localhost:8888/v0"),
		ServiceName:       "gateway",
		SkipSuccesses:     true,
		StdoutErrors:      true,
		Timeout:           10 * time.Second,
	}))

	apiAuth.Use(middleware.NewAuthMiddleware(&middleware.Config{
		AccountServiceURL: accountService,
	}))

	// routes
	accountProxy := httputil.NewSingleHostReverseProxy(proxy.MustParse(accountService))
	formProxy := httputil.NewSingleHostReverseProxy(proxy.MustParse(formService))
	mailtoProxy := httputil.NewSingleHostReverseProxy(proxy.MustParse((mailtoService)))

	// Session/Account
	api.Handle("/signup", accountProxy)
	api.Handle("/login", accountProxy)
	apiAuth.Handle("/logout", accountProxy)
	apiAuth.Handle("/getsession", accountProxy)

	// Forms/Management
	api.Handle("/form/{formID:[0-9]+}", formProxy)
	apiAuth.Handle("/forms", formProxy)
	apiAuth.Handle("/forms/{formID:[0-9]+}", formProxy)
	apiAuth.Handle("/tags", formProxy)
	apiAuth.Handle("/responses", formProxy)
	apiAuth.Handle("/responses/{responseID:[0-9]+}", formProxy)
	apiAuth.Handle("/responses/{responseID:[0-9]+}/tags", formProxy)

	// Mailto
	api.Handle("/mailto", mailtoProxy)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "Hello world!")
	})

	// Start Server
	addr := cfg.GetOrFallback("ADDR", ":80")
	log.Printf("Server %s running on %s", APIVersion, addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
