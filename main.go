package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/seanpburke/graphql-api-demo/pkg/schema"
)

func main() {
	http.HandleFunc("/", func(_ http.ResponseWriter, _ *http.Request) {}) // Handle health checks from the load balancer
	http.Handle("/query", handlers.LoggingHandler(os.Stdout, &relay.Handler{Schema: schema.Schema}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
