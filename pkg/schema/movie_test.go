package schema

import "testing"

func TestGetMovie(t *testing.T) {
	mov, err := GetMovie(2013, "Rush")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Movie: %v", mov)
	if mov.Year != 2013 {
		t.Errorf("Expecting year 2013, got %d", mov.Year)
	}
}
