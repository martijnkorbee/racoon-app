//go:build integration

//run tests with this command: go test . --tags integration -count=1

package data

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "secret"
	dbName   = "GoRacoon_TEST"
	port     = "5435"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
)

var dummyUser = User{
	FirstName: "Some",
	LastName:  "Guy",
	Email:     "me@here.com",
	Active:    1,
	Password:  "password",
}

var models Models
var testDB *sql.DB
var resource *dockertest.Resource
var pool *dockertest.Pool

func TestMain(m *testing.M) {
	os.Setenv("DATABASE_TYPE", "postgres")

	// connect to docker
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	pool = p

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13.4",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	// run postgres in docker
	resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	// wait for postgres to be alive
	if err = pool.Retry(func() (err error) {
		testDB, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			return err
		}
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to docker: %s", err)
	}

	// create postgres tables
	err = createTables(testDB)
	if err != nil {
		log.Fatalf("error creating tables: %s", err)
	}

	// create models
	models = New(testDB)

	code := m.Run()

	if err = pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func createTables(db *sql.DB) error {
	// statement
	stmt := `
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
	NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

drop table if exists users cascade;

CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	first_name character varying(255) NOT NULL,
	last_name character varying(255) NOT NULL,
	user_active integer NOT NULL DEFAULT 0,
	email character varying(255) NOT NULL UNIQUE,
	password character varying(60) NOT NULL,
	created_at timestamp without time zone NOT NULL DEFAULT now(),
	updated_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TRIGGER set_timestamp
	BEFORE UPDATE ON users
	FOR EACH ROW
	EXECUTE PROCEDURE trigger_set_timestamp();

drop table if exists remember_tokens;

CREATE TABLE remember_tokens (
	id SERIAL PRIMARY KEY,
	user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
	remember_token character varying(100) NOT NULL,
	created_at timestamp without time zone NOT NULL DEFAULT now(),
	updated_at timestamp without time zone NOT NULL DEFAULT now()
);

CREATE TRIGGER set_timestamp
	BEFORE UPDATE ON remember_tokens
	FOR EACH ROW
	EXECUTE PROCEDURE trigger_set_timestamp();

drop table if exists tokens;

CREATE TABLE tokens (
	id SERIAL PRIMARY KEY,
	user_id integer NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
	first_name character varying(255) NOT NULL,
	email character varying(255) NOT NULL,
	token character varying(255) NOT NULL,
	token_hash bytea NOT NULL,
	created_at timestamp without time zone NOT NULL DEFAULT now(),
	updated_at timestamp without time zone NOT NULL DEFAULT now(),
	expiry timestamp without time zone NOT NULL
);

CREATE TRIGGER set_timestamp
	BEFORE UPDATE ON tokens
	FOR EACH ROW
	EXECUTE PROCEDURE trigger_set_timestamp();
`

	// execute statement
	_, err := db.Exec(stmt)
	if err != nil {
		return err
	}

	return nil
}

func TestUser_Table(t *testing.T) {
	s := models.Users.Table()
	if s != "users" {
		t.Error("wrong table name returned: ", s)
	}
}

func TestUser_AddUser(t *testing.T) {
	id, err := models.Users.AddUser(dummyUser)
	if err != nil {
		t.Error("failed to insert user: ", err)
	}

	if id == 0 {
		t.Error("0 returned as id after adding user")
	}
}

func TestUser_GetAll(t *testing.T) {
	_, err := models.Users.GetAll()
	if err != nil {
		t.Error("failed to get users: ", err)
	}
}

func TestUser_GetUserByID(t *testing.T) {
	u, err := models.Users.GetUserByID(1)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	if u.ID == 0 {
		t.Error("ID of returned user is 0: ", err)
	}
}

func TestUser_GetUserByEmail(t *testing.T) {
	u, err := models.Users.GetUserByEmail(dummyUser.Email)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	if u.ID == 0 {
		t.Error("ID of returned user is 0: ", err)
	}
}

func TestUser_UpdateUser(t *testing.T) {
	u, err := models.Users.GetUserByID(1)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	u.LastName = "Smith"
	err = u.UpdateUser(*u)
	if err != nil {
		t.Error("failed to update user: ", err)
	}

	u, err = models.Users.GetUserByID(1)
	if err != nil {
		t.Error("failed to get updated user: ", err)
	}

	if u.LastName != "Smith" {
		t.Error("last name not updated in database")
	}
}

func TestUser_AuthenticateUser(t *testing.T) {
	u, err := models.Users.GetUserByID(1)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	matches, err := u.AuthenticateUser(dummyUser.Password)
	if err != nil {
		t.Error("error checking authenticate: ", err)
	}

	if !matches {
		t.Error("password doest not match when it should")
	}

	matches, err = u.AuthenticateUser("wrongpassword")
	if err != nil {
		t.Error("error checking authenticate: ", err)
	}

	if matches {
		t.Error("password matches when it should")
	}
}

func TestUser_ResetPassword(t *testing.T) {
	err := models.Users.ResetPassword(1, "new_password")
	if err != nil {
		t.Error("error resetting password: ", err)
	}

	err = models.Users.ResetPassword(2, "new_password")
	if err == nil {
		t.Error("did not get error when resetting password for non existent user")
	}
}

func TestUser_DeleteUser(t *testing.T) {
	err := models.Users.DeleteUser(1)
	if err != nil {
		t.Error("failed to delete user: ", err)
	}

	_, err = models.Users.GetUserByID(1)
	if err == nil {
		t.Error("retrieved user that was supposed to be deleted")
	}

	err = models.Users.DeleteUser(2)
	if err == nil {
		t.Error("deleted user that did not exist")
	}
}

func TestToken_Table(t *testing.T) {
	s := models.Tokens.Table()
	if s != "tokens" {
		t.Error("wrong table name returned for tokens")
	}
}

func TestToken_GenerateToken(t *testing.T) {
	id, err := models.Users.AddUser(dummyUser)
	if err != nil {
		t.Error("error inserting user: ", err)
	}

	_, err = models.Tokens.GenerateToken(id, time.Hour*24*365)
	if err != nil {
		t.Error("error generating token: ", err)
	}
}

func TestToken_InsertToken(t *testing.T) {
	u, err := models.Users.GetUserByEmail(dummyUser.Email)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	token, err := models.Tokens.GenerateToken(u.ID, time.Hour*24*365)
	if err != nil {
		t.Error("error generating token: ", err)
	}

	err = models.Tokens.InsertToken(*token, *u)
	if err != nil {
		t.Error("error inserting token: ", err)
	}
}

func TestToken_GetUserForToken(t *testing.T) {
	_, err := models.Tokens.GetUserForToken("abc")
	if err == nil {
		t.Error("no error when getting user with invalid token")
	}

	u, err := models.Users.GetUserByEmail(dummyUser.Email)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	_, err = models.Tokens.GetUserForToken(u.Token.PlainText)
	if err != nil {
		t.Error("failed to get user with valid token: ", err)
	}
}

func TestToken_GetTokensForUser(t *testing.T) {
	tokens, err := models.Tokens.GetTokensForUser(1)
	if err != nil {
		t.Error(err)
	}

	if len(tokens) > 0 {
		t.Error("tokens returned for non existing user")
	}
}

func TestToken_GetTokenByID(t *testing.T) {
	u, err := models.Users.GetUserByEmail(dummyUser.Email)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	_, err = models.Tokens.GetTokenByID(u.Token.ID)
	if err != nil {
		t.Error("error getting token by ID: ", err)
	}
}

func TestToken_GetTokenByEmail(t *testing.T) {
	u, err := models.Users.GetUserByEmail(dummyUser.Email)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	_, err = models.Tokens.GetTokenByEmail(u.Token.Email)
	if err != nil {
		t.Error("error getting token by email: ", err)
	}
}

func TestToken_GetTokenByToken(t *testing.T) {
	u, err := models.Users.GetUserByEmail(dummyUser.Email)
	if err != nil {
		t.Error("failed to get user: ", err)
	}

	_, err = models.Tokens.GetTokenByToken(u.Token.PlainText)
	if err != nil {
		t.Error("error getting token by token: ", err)
	}

	_, err = models.Tokens.GetTokenByToken("abc")
	if err == nil {
		t.Error("no error getting non existing token by token: ", err)
	}
}

var testToken_authData = []struct {
	testName    string
	token       string
	email       string
	errExpected bool
	errMessage  string
}{
	{"invalid_token", "abcdefghijklmnopqrstuvwxyz", "not@important.com", true, "invalid token accepted as valid"},
	{"invalid_token_length", "abcdefghijklmnopqrstuvwxy", "not@important.com", true, "token of wrong length accepted as valid"},
	{"no_user", "abcdefghijklmnopqrstuvwxyz", "invalid@user.com", true, "no user, but token accepted as valid"},
	{"valid_token", "", dummyUser.Email, false, "valid token reported as invalid"},
}

func TestToken_AuthenticateToken(t *testing.T) {
	for _, tt := range testToken_authData {
		var token string

		// get a user and token if needed
		if tt.email == dummyUser.Email {
			user, err := models.Users.GetUserByEmail(tt.email)
			if err != nil {
				t.Error("failed to get user: ", err)
			}
			token = user.Token.PlainText
		} else {
			token = tt.token
		}

		req, _ := http.NewRequest("GET", "/", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		_, err := models.Tokens.AuthenticateToken(req)
		if tt.errExpected && err == nil {
			t.Errorf("%s: %s", tt.testName, tt.errMessage)
		} else if !tt.errExpected && err != nil {
			t.Errorf("%s: %s - %s", tt.testName, tt.errMessage, err)
		} else {
			t.Logf("passed %s", tt.testName)
		}
	}
}

func TestToken_DeleteTokenByToken(t *testing.T) {
	user, err := models.Users.GetUserByEmail(dummyUser.Email)
	if err != nil {
		t.Error("error getting the user: ", err)
	}

	err = models.Tokens.DeleteTokenByToken(user.Token.PlainText)
	if err != nil {
		t.Error("error deleting token: ", err)
	}
}

func TestToken_ExpiredToken(t *testing.T) {
	user, err := models.Users.GetUserByEmail(dummyUser.Email)
	if err != nil {
		t.Error("error getting the user: ", err)
	}

	// create and insert a new expired token
	token, err := models.Tokens.GenerateToken(user.ID, -time.Hour*1)
	if err != nil {
		t.Error("error creating expired token: ", err)
	}

	err = models.Tokens.InsertToken(*token, *user)
	if err != nil {
		t.Error("error inserting expired token: ", err)
	}

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Authorization", "Bearer "+token.PlainText)

	_, err = models.Tokens.AuthenticateToken(req)
	if err == nil {
		t.Error("failed to catch expired token")
	}

}

func TestToken_BadHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	_, err := models.Tokens.AuthenticateToken(req)
	if err == nil {
		t.Error("failed to catch missing auth header")
	}

	req.Header.Add("Authorization", "abc")
	_, err = models.Tokens.AuthenticateToken(req)
	if err == nil {
		t.Error("failed to catch bad auth header")
	}
}

func TestToken_DeleteTokenByID(t *testing.T) {
	user, err := models.Users.GetUserByEmail(dummyUser.Email)
	if err != nil {
		t.Error("error getting the user: ", err)
	}

	err = models.Tokens.DeleteTokenByID(user.Token.ID)
	if err != nil {
		t.Error("error deleting token: ", err)
	}
}

func TestToken_ValidateToken(t *testing.T) {
	user, err := models.Users.GetUserByEmail(dummyUser.Email)
	if err != nil {
		t.Error("error getting user: ", err)
	}

	token, err := models.Tokens.GenerateToken(user.ID, 24*time.Hour)
	if err != nil {
		t.Error("error creating token: ", err)
	}

	err = models.Tokens.InsertToken(*token, *user)
	if err != nil {
		t.Error("error inserting token: ", err)
	}

	ok, err := models.Tokens.ValidateToken(token.PlainText)
	if err != nil {
		t.Error("error calling validate token: ", err)
	}

	if !ok {
		t.Error("valid token reported as invalid")
	}

	ok, err = models.Tokens.ValidateToken("invalid_token")
	if ok {
		t.Error("invalid token reported as valid")
	}
}
