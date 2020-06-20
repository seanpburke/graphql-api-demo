package schema

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/sburke-at-ziprecruiter/graphql-api-demo/pkg/table"
)

type Contact struct {
	FirstName string
	LastName  string
	Address   string
	City      string
	State     string
	Zip       string
}

type Customer struct {
	PK         string
	SK         string
	GSI2PK     string
	Phone      string // Customer's unique key.
	StorePhone string // Foreign key to customer's store.
	Contact    Contact
}

// CustomerRental represents an occasion when the customer rented movies.
type CustomerRental struct {
	PK    string
	SK    string
	Phone string
	Date  time.Time
}

func (c *Customer) MakePK() string {
	return fmt.Sprintf("CUS#%s", c.Phone)
}

func (c *Customer) MakeSK() string {
	return "CONTACT"
}

func (c *Customer) MakeGSI2PK() string {
	return fmt.Sprintf("STO#%s", c.StorePhone)
}

func (c *Customer) Init() {
	c.PK = c.MakePK()
	c.SK = c.MakeSK()
	c.GSI2PK = c.MakeGSI2PK()
	return
}

// Assure the Customer satisfies table.Item
var _ table.Item = &Customer{}

func CustomerFromJSON(jsonSrc string) (cus *Customer, err error) {
	if err = json.Unmarshal([]byte(jsonSrc), &cus); err != nil {
		return nil, fmt.Errorf("customer.FromJSON json.Unmarshal failed, %w", err)
	}
	return cus, nil
}

func GetCustomer(phone string) (cus Customer, err error) {
	cus.Phone = phone
	return cus, table.GetItem(&cus)
}

func (c Customer) Put() error {
	return table.PutItem(&c)
}

func (c Customer) Store() (Store, error) {
	return GetStore(c.StorePhone)
}

func (c Customer) PutRental() (CustomerRental, error) {
	r := CustomerRental{
		Phone: c.Phone,
		Date:  time.Now(),
	}
	return r, table.PutItem(&r)
}

func (cr *CustomerRental) MakePK() string {
	return fmt.Sprintf("CUS#%s", cr.Phone)
}

func (cr *CustomerRental) MakeSK() string {
	return fmt.Sprintf("CUS#%s#%s", cr.Phone, cr.Date)
}

func (cr *CustomerRental) Init() {
	cr.PK = cr.MakePK()
	cr.SK = cr.MakeSK()
	return
}

// Assure the CustomerRental satisfies table.Item
var _ table.Item = &CustomerRental{}
