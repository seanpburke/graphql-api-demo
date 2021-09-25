package cmd

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/seanpburke/graphql-api-demo/pkg/schema"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "serve",
		Short: "Serve GraphQL queries via HTTP",
		Args:  cobra.MaximumNArgs(1),
		Run:   Serve,
	})
}

func Serve(cmd *cobra.Command, args []string) {
	listen := viper.GetString("listen")
	if len(args) > 0 {
		listen = args[0]
	}
	log.Println("Serving HTTP on", listen)
	http.HandleFunc("/", func(_ http.ResponseWriter, _ *http.Request) {}) // Handle health checks
	http.Handle("/query", handlers.LoggingHandler(os.Stdout, &relay.Handler{Schema: schema.Schema}))
	log.Fatal(http.ListenAndServe(listen, nil))
}
