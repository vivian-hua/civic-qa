package main

import (
	// standard
	"log"
	"net/http"

	// 3rd party
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	// common
	"github.com/vivian-hua/civic-qa/services/common/environment"

	// internal
	"github.com/vivian-hua/civic-qa/services/logAggregator/internal/context"
	"github.com/vivian-hua/civic-qa/services/logAggregator/internal/repository"
)

const (
	// VersionBase is the API route base
	VersionBase = "/v0"
	// APIVersion is the API semantic version
	APIVersion = "v0.0.0"
)

var (
	// Environment
	addr   = environment.GetEnvOrFallback("ADDR", ":8888")
	dbPath = environment.GetEnvOrFallback("DBPATH", "logs.db")
)

func main() {

	// Routers
	router := mux.NewRouter()
	api := router.PathPrefix(VersionBase).Subrouter()

	// Handler context
	repo, err := repository.NewLogRepository(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to create log repository: %v", err)
	}
	ctx := &context.Context{Repo: repo}

	// Routes
	api.HandleFunc("/log", ctx.HandleLog)
	api.HandleFunc("/query", ctx.HandleQuery)

	// Start Server
	log.Printf("Server %s running on %s", APIVersion, addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
