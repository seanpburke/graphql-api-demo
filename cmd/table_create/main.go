package main

import (
	"fmt"
	"os"

	"github.com/seanpburke/graphql-api-demo/pkg/config"
	"github.com/seanpburke/graphql-api-demo/pkg/table"
)

func main() {

	out, err := table.CreateTable()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Created the DynamoDB table", config.Config.Table)
	fmt.Println(out)
}
