package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	Verbose bool

	rootCmd = &cobra.Command{
		Use:   "graphql",
		Short: "Demo GraphQL application",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
)

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
