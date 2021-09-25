package cmd

import (
	"fmt"
	"log"

	"github.com/seanpburke/graphql-api-demo/pkg/schema"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "movie-rent",
		Short: "Record a movie rental in DynamoDB",
		Args:  cobra.NoArgs,
		Run:   MovieRent,
	})
}

func MovieRent(cmd *cobra.Command, args []string) {

	// Get my customer
	customerPhone := "828-234-1717"
	cus, err := schema.GetCustomer(customerPhone)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Customer:", cus)

	sto, err := cus.Store()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Store:", sto)

	rental, err := cus.PutRental()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Rental:", rental)

	mov, err := schema.GetMovie(2013, "Rush")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Movie:", mov)

	movren, err := mov.PutRental(rental.Phone, rental.Date)
	if err != nil {
		log.Fatal(err.Error())
	}
	if Verbose {
		fmt.Printf("Successfully added rental %s (%d) to %s on %s\n", movren.Title, movren.Year, movren.Phone, movren.Date)
	}
}
