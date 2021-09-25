package cmd

import (
	"log"

	"github.com/seanpburke/graphql-api-demo/pkg/schema"
	"github.com/spf13/cobra"
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

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "customers-load",
		Short: "Load customers into DynamoDB",
		Args:  cobra.NoArgs,
		Run:   CustomersLoad,
	})
}

func CustomersLoad(cmd *cobra.Command, args []string) {

	// Load customers into the DynamoDB table.
	for _, cus := range customers {
		if err := cus.Put(); err != nil {
			log.Fatal(err.Error())
		}
		// Load the customer's store
		sto, err := cus.Store()
		if err != nil {
			log.Fatal(err.Error())
		}
		// Get the store's customers
		customers, err := sto.Customers()
		if err != nil {
			log.Fatal(err.Error())
		}
		if len(customers) == 0 {
			log.Fatal("Expecting at least one customer for store", sto.Name)
		}
		c := customers[0]
		if Verbose {
			log.Printf("Successfully added %s %s (%s) to %s\n", c.Contact.FirstName, c.Contact.LastName, c.Phone, sto.Name)
		}
	}
}
