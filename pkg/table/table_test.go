package table

import "testing"

type Thing struct {
	PK    string
	SK    string
	Phone string
}

func (t Thing) MakePK() string {
	return t.PK
}
func (t Thing) MakeSK() string {
	return t.SK
}
func (t Thing) Init() {
}

func TestGetItem(t *testing.T) {
	phone := "828-234-1717"
	thing := Thing{
		PK: "CUS#" + phone,
		SK: "CONTACT",
	}
	if err := GetItem(&thing); err != nil {
		t.Error(err)
		return
	}
	// t.Logf("Thing: %v", thing)
	if thing.Phone != phone {
		t.Errorf("Expecting phone %q, got %q", phone, thing.Phone)
	}
}
