package schema

import (
	"github.com/graph-gophers/graphql-go"
)

//go:generate ../../scripts/schema_graphql.sh

type RootResolver struct{}

func (r *RootResolver) Customer(args struct{ Phone string }) (Customer, error) {
	return GetCustomer(args.Phone)
}

func (r *RootResolver) Movie(args struct {
	Year  int32
	Title string
}) (Movie, error) {
	return GetMovie(args.Year, args.Title)
}

func (r *RootResolver) Store(args struct{ Phone string }) (Store, error) {
	return GetStore(args.Phone)
}

func (s Store) Movies(args struct {
	Year  int32
	Title string
}) ([]StoreMovie, error) {
	return s.GetMovies(args.Year, args.Title)
}

var (
	// We can pass an option to the schema so we don’t need to
	// write a method to access each type’s field:
	opts   = []graphql.SchemaOpt{graphql.UseFieldResolvers()}
	Schema = graphql.MustParseSchema(schemaString, &RootResolver{}, opts...)
)
