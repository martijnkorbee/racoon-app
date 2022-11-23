package data

import (
	"database/sql"
	"fmt"

	db2 "github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

var db *sql.DB
var upper db2.Session

type Models struct {
	// any models inserted here (and in the New function)
	// are easily accessible throughout the entire application
	Users  User
	Tokens Token
}

func New(dbPool *sql.DB) Models {
	db = dbPool

	upper, _ = postgresql.New(db)

	return Models{
		Users:  User{},
		Tokens: Token{},
	}
}

func getInsertID(i db2.ID) int {
	idType := fmt.Sprintf("%T", i)
	if idType == "int64" {
		return int(i.(int64))
	}

	return i.(int)
}
