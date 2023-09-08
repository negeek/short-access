package user

import(
	//"fmt"
	"context"
	"github.com/negeek/short-access/utils"
	"github.com/negeek/short-access/db"
	"github.com/jackc/pgx/v4"
)

// create a new user
func (u *User) Create() error {
	// set the date fields
	terr:=utils.Time(u,true)
	if terr != nil {
		return terr
	}
	// insert user detail into users table in db
	query:="INSERT INTO users (id, password, email, date_created, date_updated) VALUES ($1, $2, $3, $4, $5)"
	_, ierr := db.PostgreSQLDB.Exec(context.Background(),query,u.Id, u.Password, u.Email, u.DateCreated, u.DateUpdated)
	if ierr != nil {
		return ierr
	}
	return nil
}

// check if email exists
func (u *User) EmailExists() bool{
	var emailExists bool
	query:="SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)"
	err:= db.PostgreSQLDB.QueryRow(context.Background(),query,u.Email).Scan(&emailExists)
	if err !=nil{
		return false
	}
	return emailExists
}

// find user by email. email is also unique like id
func (u *User) FindByEmail() (error, bool) {
	query:="SELECT id, email, date_created, date_updated FROM users WHERE email = $1"
	err := db.PostgreSQLDB.QueryRow(context.Background(),query,u.Email).Scan(&u.Id, &u.Email, &u.DateCreated, &u.DateUpdated)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, false
		}
		return err,false
	}
	return nil, true
}

func (u *User) Authenticate() (error, bool) {
	query:="SELECT id, email, date_created, date_updated FROM users WHERE email = $1 and password = $2"
	err := db.PostgreSQLDB.QueryRow(context.Background(),query,u.Email,u.Password).Scan(&u.Id, &u.Email, &u.DateCreated, &u.DateUpdated)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, false
		}
		return err,false
	}
	return nil, true
}

func (u *User) Delete() error {
	query:="DELETE FROM users WHERE email = $1"
	_, err := db.PostgreSQLDB.Exec(context.Background(), query, u.Email)
	if err != nil {
		return err
	}
	return nil
}

