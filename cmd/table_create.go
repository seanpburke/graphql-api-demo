package cmd

import (
	"fmt"
	"log"

	"github.com/seanpburke/graphql-api-demo/pkg/config"
	"github.com/seanpburke/graphql-api-demo/pkg/table"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "table-create",
		Short: "Create the DynamoDB table",
		Args:  cobra.NoArgs,
		Run:   TableCreate,
	})
}

func TableCreate(cmd *cobra.Command, args []string) {

	out, err := table.CreateTable()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Created the DynamoDB table", config.Config.Table)
	if Verbose {
		fmt.Println(out)
	}
}
