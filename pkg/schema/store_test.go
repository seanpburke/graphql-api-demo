package schema

import "testing"

func TestGetStore(t *testing.T) {
	phone := "828-555-1249"

	sto, err := GetStore(phone)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Store: %v", sto)
	if sto.Phone != phone {
		t.Errorf("Expecting phone %q, got %q", phone, sto.Phone)
	}

	cus, err := sto.Customers()
	if err != nil {
		t.Error(err)
		return
	}
	if cus[0].StorePhone != phone {
		t.Errorf("Expecting customer StorePhone %q, got %q", phone, cus[0].StorePhone)
	}

	sm, err := sto.Movies(struct {
		Year  int32
		Title string
	}{Year: 2013, Title: "Rush"})
	if err != nil {
		t.Error(err)
		return
	}
	if sm[0].Title != "Rush" {
		t.Errorf("Expecting movie title %q, got %q", "Rush", sm[0].Title)
	}
}
