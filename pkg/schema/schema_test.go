package schema

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

// These tests are based on the helpful examples in
// https://github.com/codex-zaydek/graphql-go-walkthrough

func TestCustomer(t *testing.T) {

	ctx := context.Background()
	query := `
	query Customer($phone: String!) {
		customer(phone: $phone) {
			phone
			storephone
			contact {
				firstname lastname address city state zip
			}
			store { name location { address } }
		}
	}`
	phone := "815-717-3861"
	param := map[string]interface{}{
		"phone": phone,
	}
	resp := Schema.Exec(ctx, query, "Customer", param)
	if !strings.Contains(string(resp.Data), phone) {
		t.Errorf("Customer phone %q not found", phone)
	}
	json, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		t.Error(err)
	}
	t.Logf("Customer: %v", string(json))
}

func TestMovie(t *testing.T) {

	ctx := context.Background()
	query := `
	query Movie($year: Int!, $title: String!) {
		movie(year: $year, title: $title) {
			year title info { directors rating genres plot rank actors }
		}
	}`
	year := 2013
	title := "Rush"
	param := map[string]interface{}{
		"year":  year,
		"title": title,
	}
	resp := Schema.Exec(ctx, query, "Movie", param)
	if !strings.Contains(string(resp.Data), title) {
		t.Errorf("Movie %q(%d) not found", title, year)
	}
	json, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		t.Error(err)
	}
	t.Logf("Movie: %v", string(json))
}

func TestStore(t *testing.T) {

	ctx := context.Background()
	query := `
	query Store($phone: String!) {
		store(phone: $phone) {
			phone name
			location  { address city state zip }
			customers { phone contact { firstname lastname } }
		}
	}`
	phone := "828-555-1249"
	param := map[string]interface{}{
		"phone": phone,
	}
	resp := Schema.Exec(ctx, query, "Store", param)
	if !strings.Contains(string(resp.Data), phone) {
		t.Errorf("Store phone %q not found", phone)
	}
	json, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		t.Error(err)
	}
	t.Logf("Store: %v", string(json))
}

func TestInventory(t *testing.T) {

	ctx := context.Background()
	query := `
	query Inventory($phone: String!, $year: Int!, $title: String!) {
		store(phone: $phone) {
			phone name
			movies(year: $year, title: $title) { title, count }
		}
	}`
	phone := "828-555-1249"
	param := map[string]interface{}{
		"phone": phone,
		"year":  2014,
		"title": "",
	}
	resp := Schema.Exec(ctx, query, "Inventory", param)
	if !strings.Contains(string(resp.Data), phone) {
		t.Errorf("Store phone %q not found", phone)
	}
	json, err := json.MarshalIndent(resp, "", "\t")
	if err != nil {
		t.Error(err)
	}
	t.Logf("Inventory: %v", string(json))
}
