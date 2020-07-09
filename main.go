package main

import (
	"log"
	"net/http"

	"github.com/graph-gophers/graphql-go/relay"
	"github.com/sburke-at-ziprecruiter/graphql-api-demo/pkg/schema"
)

func main() {
	http.HandleFunc("/", func(_ http.ResponseWriter, _ *http.Request) {}) // Handle health checks from the load balancer
	http.Handle("/query", &relay.Handler{Schema: schema.Schema})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
