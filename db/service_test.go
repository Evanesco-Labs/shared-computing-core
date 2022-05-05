package db

import (
	"testing"
)

func TestDBService_GetUserSpaceship(t *testing.T) {
	_, err := Init("", "","")
	if err != nil {
		t.Fatal(err)
	}
}
