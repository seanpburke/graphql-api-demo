package main

import (
	"fmt"
	"os"

	"github.com/sburke-at-ziprecruiter/graphql-api-demo/pkg/schema"
)

var customers = []schema.Customer{
	{
		Phone:      "828-234-1717",
		StorePhone: "828-555-1249", // Brilliant Video
		Contact: schema.Contact{
			FirstName: "Alex",
			LastName:  "Able",
			Address:   "303 British Street",
			City:      "Lordston",
			State:     "Michigan",
			Zip:       "28202",
		},
	},
	{
		Phone:      "414-232-1858",
		StorePhone: "310-555-8800", // Dazzling Video
		Contact: schema.Contact{
			FirstName: "Betty",
			LastName:  "Bialoski",
			Address:   "762 Nato Street",
			City:      "Chaffing",
			State:     "Montana",
			Zip:       "58201",
		},
	},
	{
		Phone:      "815-717-3861",
		StorePhone: "415-555-0117", // Sizzling Video
		Contact: schema.Contact{
			FirstName: "Charlie",
			LastName:  "Chalice",
			Address:   "28 Floze Ave",
			City:      "Avalon",
			State:     "Tenessee",
			Zip:       "40385",
		},
	},
}

func main() {

	// Load customers into the DynamoDB table.
	for _, cus := range customers {
		if err := cus.Put(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		// Load the customer's store
		sto, err := cus.Store()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		// Get the store's customers
		customers, err := sto.Customers()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if len(customers) != 1 {
			fmt.Println("Expecting one customer for store", sto.Name)
			os.Exit(1)
		}
		c := customers[0]
		fmt.Printf("Successfully added %s %s (%s) to %s\n", c.Contact.FirstName, c.Contact.LastName, c.Phone, sto.Name)
	}
}
