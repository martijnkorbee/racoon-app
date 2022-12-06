package data

import (
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

func TestNew(t *testing.T) {
	fakeDB, _, _ := sqlmock.New()
	defer fakeDB.Close()

	upper, _ := postgresql.New(fakeDB)

	_ = os.Setenv("DATABASE_TYPE", "postgres")
	_ = os.Setenv("UPPER_DB_LOG", "ERROR")
	m := New(upper)

	if fmt.Sprintf("%T", m) != "data.Models" {
		t.Error("wrong type", fmt.Sprintf("%T", m))
	}
}

func TestGetInsertID(t *testing.T) {
	var id db.ID
	id = int64(1)

	returnedID := getInsertID(id)
	if fmt.Sprintf("%T", returnedID) != "int" {
		t.Error("wrong type returned")
	}

	id = 1
	returnedID = getInsertID(id)
	if fmt.Sprintf("%T", returnedID) != "int" {
		t.Error("wrong type returned")
	}
}
