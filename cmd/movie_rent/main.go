package main

import (
	"fmt"
	"os"

	"github.com/sburke-at-ziprecruiter/graphql-api-demo/pkg/schema"
)

func main() {

	// Get my customer
	customerPhone := "828-234-1717"
	cus, err := schema.GetCustomer(customerPhone)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Customer:", cus)

	sto, err := cus.Store()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Store:", sto)

	rental, err := cus.PutRental()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Rental:", rental)

	mov, err := schema.GetMovie(2013, "Rush")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Movie:", mov)

	movren, err := mov.PutRental(rental.Phone, rental.Date)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("Successfully added rental %s (%d) to %s on %s\n", movren.Title, movren.Year, movren.Phone, movren.Date)
}
