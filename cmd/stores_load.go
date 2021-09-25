package cmd

import (
	"log"

	"github.com/seanpburke/graphql-api-demo/pkg/schema"
	"github.com/spf13/cobra"
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

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "stores-load",
		Short: "Load stores into DynamoDB",
		Args:  cobra.NoArgs,
		Run:   StoresLoad,
	})
}

func StoresLoad(cmd *cobra.Command, args []string) {

	// Scan all of the movies
	movies, err := schema.ScanMovies()
	if err != nil {
		log.Fatal(err.Error())
	}

	// Load stores into DynamoDB
	for _, sto := range stores {

		if err := sto.Put(); err != nil {
			log.Fatal(err.Error())
		}
		if Verbose {
			log.Printf("Successfully added %q (%s)\n", sto.Name, sto.PK)
		}

		// Put 3 copies of each movie in this store's inventory.
		for _, mov := range movies {
			if err := sto.PutMovie(mov.Year, mov.Title, 3); err != nil {
				log.Fatal(err.Error())
			}
		}

		// Fetch this store's entire inventory.
		inventory, err := sto.GetMovies(0, "")
		if err != nil {
			log.Fatal(err.Error())
		}
		if Verbose {
			for _, inv := range inventory {
				log.Printf("Successfully added %q(%d)[%d] to inventory of %q\n", inv.Title, inv.Year, inv.Count, sto.Name)
			}
		}
	}
}
