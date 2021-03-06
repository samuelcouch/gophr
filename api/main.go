package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gophr-pm/gophr/lib"
	"github.com/gophr-pm/gophr/lib/datadog"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize the API.
	config, client := lib.Init()

	// Ensure that the client is closed eventually.
	defer client.Close()

	// Initialize datadog client.
	dataDogClient, err := datadog.NewClient(config, "api.")
	if err != nil {
		log.Println(err)
	}

	// Register all of the routes.
	r := mux.NewRouter()
	r.HandleFunc("/status", StatusHandler()).Methods("GET")
	r.HandleFunc(fmt.Sprintf(
		"/blob/{%s}/{%s}/{%s}/{%s}",
		urlVarAuthor,
		urlVarRepo,
		urlVarSHA,
		urlVarPath),
		BlobHandler(dataDogClient)).Methods("GET")
	r.HandleFunc(
		"/packages/new",
		GetNewPackagesHandler(client, dataDogClient)).Methods("GET")
	r.HandleFunc(
		"/packages/search",
		SearchPackagesHandler(client, dataDogClient)).Methods("GET")
	r.HandleFunc(
		"/packages/trending",
		GetTrendingPackagesHandler(client, dataDogClient)).Methods("GET")
	r.HandleFunc(fmt.Sprintf(
		"/packages/top/{%s}/{%s}",
		urlVarLimit,
		urlVarTimeSplit),
		GetTopPackagesHandler(client, dataDogClient)).Methods("GET")
	r.HandleFunc(fmt.Sprintf(
		"/packages/{%s}/{%s}",
		urlVarAuthor,
		urlVarRepo),
		GetPackageHandler(client, dataDogClient)).Methods("GET")

	// Start serving.
	log.Printf("Servicing HTTP requests on port %d.\n", config.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r)
}
