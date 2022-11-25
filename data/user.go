package data

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	up "github.com/upper/db/v4"
)

type User struct {
	ID        int       `db:"id,omitempty"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Active    int       `db:"user_active"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Token     Token     `db:"-"`
}

func (u *User) Table() string {
	return "users"
}

// AddUser adds a user
func (u *User) AddUser(user User) (ID int, err error) {
	// create password hash
	newHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Password = string(newHash)

	collection := upper.Collection(u.Table())
	res, err := collection.Insert(user)
	if err != nil {
		return 0, err
	}

	id := getInsertID(res.ID())

	return id, nil
}

// GetAll gets all users
func (u *User) GetAll() ([]*User, error) {
	var all []*User
	collection := upper.Collection(u.Table())

	res := collection.Find().OrderBy("last_name")
	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// GetUserByID gets user by ID
func (u *User) GetUserByID(ID int) (*User, error) {
	var user User
	collection := upper.Collection(u.Table())

	res := collection.Find(up.Cond{"id =": ID})
	err := res.One(&user)
	if err != nil {
		return nil, err
	}

	var token Token
	collection = upper.Collection(token.Table())

	res = collection.Find(up.Cond{"user_id =": user.ID, "expiry >": time.Now()}).OrderBy("created_at desc")
	err = res.One(&token)
	if err != nil {
		if err != up.ErrNilRecord && err != up.ErrNoMoreRows {
			return nil, err
		}
	}

	user.Token = token

	return &user, nil
}

// GetUserByEmail gets user by email
func (u *User) GetUserByEmail(email string) (*User, error) {
	var user User
	collection := upper.Collection(u.Table())

	res := collection.Find(up.Cond{"email =": email})
	err := res.One(&user)
	if err != nil {
		return nil, err
	}

	var token Token
	collection = upper.Collection(token.Table())

	res = collection.Find(up.Cond{"user_id =": user.ID, "expiry >": time.Now()}).OrderBy("created_at desc")
	err = res.One(&token)
	if err != nil {
		if err != up.ErrNilRecord && err != up.ErrNoMoreRows {
			return nil, err
		}
	}

	user.Token = token

	return &user, nil
}

// UpdateUser updates a user
func (u *User) UpdateUser(user User) error {
	user.UpdatedAt = time.Now()

	collection := upper.Collection(u.Table())
	res := collection.Find(user.ID)
	err := res.Update(&user)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes a user
func (u *User) DeleteUser(id int) error {
	collection := upper.Collection(u.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}

	return nil
}

// ResetPassword resets the users password
func (u *User) ResetPassword(id int, password string) error {
	// create new password hash
	newHash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	user, err := u.GetUserByID(id)
	if err != nil {
		return err
	}

	user.Password = string(newHash)

	err = user.UpdateUser(*user)
	if err != nil {
		return err
	}

	return nil
}

// AuthenticateUser authenticates user
func (u *User) AuthenticateUser(password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
