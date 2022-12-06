package data

import (
	"fmt"

	"github.com/upper/db/v4"
)

type Models struct {
	// any models inserted here (and in the New function)
	// are easily accessible throughout the entire application
	Users  User
	Tokens Token
}

var upper db.Session

func New(db db.Session) Models {
	upper = db

	return Models{
		Users:  User{},
		Tokens: Token{},
	}
}

func getInsertID(i db.ID) int {
	idType := fmt.Sprintf("%T", i)
	if idType == "int64" {
		return int(i.(int64))
	}

	return i.(int)
}
