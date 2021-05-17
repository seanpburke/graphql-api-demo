package schema

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/seanpburke/graphql-api-demo/pkg/config"
	"github.com/seanpburke/graphql-api-demo/pkg/table"
)

type Location struct {
	Address string
	City    string
	State   string
	Zip     string
}

type Store struct {
	PK       string
	SK       string
	Phone    string
	Name     string
	Location Location
}

// StoreMovie represents the copies of a movie in a store's inventory.
type StoreMovie struct {
	PK    string
	SK    string
	Phone string // Pkey to Store
	Year  int32  // Pkey to Movie
	Title string // Pkey to Movie
	Count int32
}

func (s *Store) MakePK() string {
	return fmt.Sprintf("STO#%s", s.Phone)
}

func (s *Store) MakeSK() string {
	return "LOCATION"
}

func (s *Store) Init() {
	s.PK = s.MakePK()
	s.SK = s.MakeSK()
	return
}

// Assure the Store satisfies table.Item
var _ table.Item = &Store{}

func StoreFromJSON(jsonSrc string) (store *Store, err error) {
	if err = json.Unmarshal([]byte(jsonSrc), &store); err != nil {
		return nil, fmt.Errorf("store.FromJSON json.Unmarshal failed, %w", err)
	}
	return store, nil
}

func GetStore(phone string) (sto Store, err error) {
	sto.Phone = phone
	return sto, table.GetItem(&sto)
}

func (s *Store) Put() error {
	return table.PutItem(s)
}

func (s Store) Customers() ([]Customer, error) {
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":k": {
				S: aws.String(s.PK),
			},
		},
		KeyConditionExpression: aws.String("GSI2PK = :k"),
		IndexName:              aws.String("GSI2"),
		TableName:              aws.String(config.Config.Table),
	}
	result, err := config.Config.DDB.Query(input)
	if err != nil {
		return nil, err
	}
	customers := make([]Customer, len(result.Items))
	for ix, item := range result.Items {
		if err := dynamodbattribute.UnmarshalMap(item, &customers[ix]); err != nil {
			return nil, err
		}
	}
	return customers, nil
}

func (s Store) PutMovie(year int32, title string, count int32) error {
	sm := StoreMovie{
		Phone: s.Phone,
		Year:  year,
		Title: title,
		Count: count,
	}
	return sm.Put()
}

// Movies returns the store's inventory.
// The year and title can be used to constrain the results,
// but because we apply this constraint to the sort key,
// the year must be specified, and if title is specied,
// any movie titles with that prefix will be returned.
// Obviously, this is only an exercise in DDB query techniques,
// and very much not a practical general-purpose search capability.
func (s *Store) GetMovies(year int32, title string) ([]StoreMovie, error) {
	// Use year and title to construct a prefix test for the sort key.
	sk := "MOV#"
	if year != 0 {
		sk = sk + fmt.Sprintf("%d", year)
		if title != "" {
			sk = sk + "#" + title
		}
	}
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(s.PK),
			},
			":sk": {
				S: aws.String(sk),
			},
		},
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK,:sk)"),
		TableName:              aws.String(config.Config.Table),
	}
	result, err := config.Config.DDB.Query(input)
	if err != nil {
		return nil, err
	}
	movies := make([]StoreMovie, len(result.Items))
	for ix, item := range result.Items {
		if err := dynamodbattribute.UnmarshalMap(item, &movies[ix]); err != nil {
			return nil, err
		}
	}
	return movies, nil
}

func (sm *StoreMovie) MakePK() string {
	return fmt.Sprintf("STO#%s", sm.Phone)
}

func (sm *StoreMovie) MakeSK() string {
	return fmt.Sprintf("MOV#%d#%s", sm.Year, sm.Title)
}

func (sm *StoreMovie) Init() {
	sm.PK = sm.MakePK()
	sm.SK = sm.MakeSK()
	return
}

// Assure the Store satisfies table.Item
var _ table.Item = &StoreMovie{}

func (sm *StoreMovie) Put() error {
	return table.PutItem(sm)
}

func (sm StoreMovie) GetMovie() (mov Movie, err error) {
	mov.Year = sm.Year
	mov.Title = sm.Title
	return mov, table.GetItem(&mov)
}
