package main

import (
	"github.com/betacraft/yaag/middleware"
	"github.com/betacraft/yaag/yaag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	handlerWrapper, err := NewHandlerWrapper()
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	yaag.Init(&yaag.Config{On: true, DocTitle: "Downsampling Misc Services", DocPath: "servicedoc.html"})
	router := mux.NewRouter().StrictSlash(true)
	// Get object
	router.HandleFunc("/metrics-downsampling/{queryid}", middleware.HandleFunc(handlerWrapper.GetMetricsDownsamplingItem)).Methods("GET", "OPTIONS")
	// preview object
	router.HandleFunc("/metrics-downsampling/previews/pending/_next", middleware.HandleFunc(handlerWrapper.GetMetricsDownsamplingForPreview)).Methods("GET", "OPTIONS")
	// query object for flink
	router.HandleFunc("/metrics-downsampling/queries/pending/_next", middleware.HandleFunc(handlerWrapper.GetMetricsDownsamplingPendingQuery)).Methods("GET", "OPTIONS")
	// update pending to deployed
	router.HandleFunc("/metrics-downsampling/previews/pending/{queryid}/update/_deployed", middleware.HandleFunc(handlerWrapper.DeployMetricsDownsamplingPreview)).
		Methods("POST", "OPTIONS")

	// get all metrics stacks
	router.HandleFunc("/metrics-stacks", middleware.HandleFunc(handlerWrapper.GetMetricsStacks)).Methods("GET", "OPTIONS")

	log.Fatal(http.ListenAndServe(":8080", handlerWrapper.WrapHandler(router)))
}
