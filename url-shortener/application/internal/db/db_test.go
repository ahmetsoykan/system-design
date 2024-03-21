package db

import (
	"testing"
)

// connects the same table defined in write.go file
// TODO: add more test cases
func TestDB(t *testing.T) {

	// Client
	d := NewDynamoDBClient("eu-west-1")

	d.DDBInit()

	// Write operation
	putItem, err := d.DDBPutItem(Item{
		ID:         "",
		LongURL:    "http://facebook.com/",
		Short:      "",
		DefaultTTL: true,
		TTL:        0,
	})
	if err != nil {
		t.Errorf("failed: got %s, want nil", err.Error())
	}

	// Read operation
	item, _ := d.GetItembyPK(putItem.ID)
	if item.LongURL != "http://facebook.com/" {
		t.Errorf("failed: got %s, want %s", err.Error(), "http://facebook.com")
	}
}
