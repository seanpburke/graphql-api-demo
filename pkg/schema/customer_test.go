package schema

import "testing"

func TestGetCustomer(t *testing.T) {
	phone := "828-234-1717"

	cus, err := GetCustomer(phone)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Customer: %v", cus)
	if cus.Phone != phone {
		t.Errorf("Expecting phone %q, got %q", phone, cus.Phone)
	}
}
