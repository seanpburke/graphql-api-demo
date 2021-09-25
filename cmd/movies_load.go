package cmd

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/seanpburke/graphql-api-demo/pkg/schema"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "movies-load",
		Short: "Load movies into DynamoDB",
		Args:  cobra.NoArgs,
		Run:   MoviesLoad,
	})
}

// readMovies reads Movies from a gzipped JSON via stdin.
// Note that these movies need to call Init() to set PK and SK.
func readMovies() ([]schema.Movie, error) {
	gunzip, err := gzip.NewReader(os.Stdin)
	if err != nil {
		return nil, err
	}
	defer gunzip.Close()

	raw, err := ioutil.ReadAll(gunzip)
	if err != nil {
		return nil, err
	}

	var movies []schema.Movie
	err = json.Unmarshal(raw, &movies)
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func MoviesLoad(cmd *cobra.Command, args []string) {

	// Load Movies via stdin and add them the DDB table.
	movies, err := readMovies()
	if err != nil {
		log.Fatal(err.Error())
	}
	for ix, movie := range movies {
		if ix > 9 {
			break
		}
		if err = movie.Put(); err != nil {
			log.Fatal(err.Error())
		}
	}

	// Scan the movies that we loaded.
	movies, err = schema.ScanMovies()
	if err != nil {
		log.Fatal(err.Error())
	}
	if Verbose {
		for _, movie := range movies {
			fmt.Printf("Successfully added %q (%d) to table.\n", movie.Title, movie.Year)
		}
	}
}
