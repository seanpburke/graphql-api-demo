package schema

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/sburke-at-ziprecruiter/graphql-api-demo/pkg/config"
	"github.com/sburke-at-ziprecruiter/graphql-api-demo/pkg/table"
)

type Info struct {
	Directors       []string
	ReleaseDate     time.Time `json:"release_date" dynamodbav:"ReleaseDate"`
	Rating          float64
	Genres          []string
	ImageURL        string `json:"image_url" dynamodbav:"ImageURL"`
	Plot            string
	Rank            int32
	RunningTimeSecs int `json:"running_time_secs" dynamodbav:"RunningTimeSecs"`
	Actors          []string
}

type Movie struct {
	PK    string
	SK    string
	Year  int32
	Title string
	Info  Info
}

// MovieRental represents a movie that was part of a CustomerRental
type MovieRental struct {
	PK         string
	SK         string
	Year       int32     // part of pkey to Movie
	Title      string    // part of pkey to Movie
	Phone      string    // Part of fkey to CustomerRental
	Date       time.Time // Part of fkey to CustomerRental
	DueDate    time.Time
	ReturnDate time.Time
}

func (m *Movie) MakePK() string {
	return fmt.Sprintf("MOV#%d#%s", m.Year, m.Title)
}

func (m *Movie) MakeSK() string {
	return "INFO"
}

func (m *Movie) Init() {
	m.PK = m.MakePK()
	m.SK = m.MakeSK()
	return
}

// Assure that Movie satisfies table.Item
var _ table.Item = &Movie{}

func MovieFromJSON(jsonSrc string) (movie *Movie, err error) {
	if err = json.Unmarshal([]byte(jsonSrc), &movie); err != nil {
		return nil, fmt.Errorf("movie.FromJSON json.Unmarshal failed, %w", err)
	}
	return movie, nil
}

func GetMovie(year int32, title string) (mov Movie, err error) {
	mov.Year = year
	mov.Title = title
	return mov, table.GetItem(&mov)
}

func (m *Movie) Put() error {
	return table.PutItem(m)
}

// Scan the table for all movie INFO items.
func ScanMovies() ([]Movie, error) {

	input := &dynamodb.ScanInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":sk": {
				S: aws.String("INFO"),
			},
		},
		FilterExpression: aws.String("SK = :sk"),
		TableName:        aws.String(config.Config.Table),
	}

	result, err := config.Config.DDB.Scan(input)
	if err != nil {
		return nil, err
	}
	movies := make([]Movie, len(result.Items))
	for ix, item := range result.Items {
		if err := dynamodbattribute.UnmarshalMap(item, &movies[ix]); err != nil {
			return nil, err
		}
	}
	return movies, nil
}

func (mr *MovieRental) MakePK() string {
	return fmt.Sprintf("MOV#%d#%s", mr.Year, mr.Title)
}

func (mr *MovieRental) MakeSK() string {
	return fmt.Sprintf("REN#%s#%s", mr.Phone, mr.Date)
}

func (mr *MovieRental) Init() {
	mr.PK = mr.MakePK()
	mr.SK = mr.MakeSK()
	return
}

// Assure that MovieRental satisfies table.Item
var _ table.Item = &MovieRental{}

func (m Movie) PutRental(phone string, date time.Time) (MovieRental, error) {
	mr := MovieRental{
		Year:    m.Year,
		Title:   m.Title,
		Phone:   phone,
		Date:    date,
		DueDate: time.Now().Add(30 * 24 * time.Hour),
	}
	return mr, table.PutItem(&mr)
}
