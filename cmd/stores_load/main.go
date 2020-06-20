package main

import (
	"fmt"
	"os"

	"github.com/sburke-at-ziprecruiter/graphql-api-demo/pkg/schema"
)

var stores = []schema.Store{
	{
		Name:  "Brilliant Video",
		Phone: "828-555-1249",
		Location: schema.Location{
			Address: "123 Main Street",
			City:    "Townville",
			State:   "Connecticut",
			Zip:     "06010",
		},
	},
	{
		Name:  "Dazzling Video",
		Phone: "310-555-8800",
		Location: schema.Location{
			Address: "777 Lucky Blvd",
			City:    "Lost Angels",
			State:   "California",
			Zip:     "90045",
		},
	},
	{
		Name:  "Sizzling Video",
		Phone: "415-555-0117",
		Location: schema.Location{
			Address: "5500 Market Street",
			City:    "San Francisco",
			State:   "California",
			Zip:     "91199",
		},
	},
}

func main() {

	// Scan all of the movies
	movies, err := schema.ScanMovies()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Load stores into DynamoDB
	for _, sto := range stores {

		if err := sto.Put(); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Printf("Successfully added %q (%s)\n", sto.Name, sto.PK)

		// Put 3 copies of each movie in this store's inventory.
		for _, mov := range movies {
			if err := sto.PutMovie(mov.Year, mov.Title, 3); err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
		}

		// Fetch this store's entire inventory.
		inventory, err := sto.GetMovies(0, "")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		for _, inv := range inventory {
			fmt.Printf("Successfully added %q(%d)[%d] to inventory of %q\n", inv.Title, inv.Year, inv.Count, sto.Name)
		}
	}
}
