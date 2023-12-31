package pgsql

import (
	"net/http"
	"strings"

	"github.com/go-pg/pg/v9"

	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo/v4"

	stems "github.com/Soapstone-Services/go-template-2024"
)

// User represents the client for user table
type User struct{}

// Custom errors
var (
	ErrAlreadyExists = echo.NewHTTPError(http.StatusInternalServerError, "Username or email already exists.")
)

// Create creates a new user on database
func (u User) Create(db orm.DB, usr stems.User) (stems.User, error) {
	var user = new(stems.User)
	err := db.Model(user).Where("lower(username) = ? or lower(email) = ? and deleted_at is null",
		strings.ToLower(usr.Username), strings.ToLower(usr.Email)).Select()
	if err == nil || err != pg.ErrNoRows {
		return stems.User{}, ErrAlreadyExists
	}

	err = db.Insert(&usr)
	return usr, err
}

// View returns single user by ID
func (u User) View(db orm.DB, id int) (stems.User, error) {
	var user stems.User
	sql := `SELECT "user".*, "role"."id" AS "role__id", "role"."access_level" AS "role__access_level", "role"."name" AS "role__name" 
	FROM "users" AS "user" LEFT JOIN "roles" AS "role" ON "role"."id" = "user"."role_id" 
	WHERE ("user"."id" = ? and deleted_at is null)`
	_, err := db.QueryOne(&user, sql, id)
	return user, err
}

// Update updates user's contact info
func (u User) Update(db orm.DB, user stems.User) error {
	_, err := db.Model(&user).WherePK().UpdateNotZero()
	return err
}

// List returns list of all users retrievable for the current user, depending on role
func (u User) List(db orm.DB, qp *stems.ListQuery, p stems.Pagination) ([]stems.User, error) {
	var users []stems.User
	q := db.Model(&users).Relation("Role").Limit(p.Limit).Offset(p.Offset).Where("deleted_at is null").Order("user.id desc")
	if qp != nil {
		q.Where(qp.Query, qp.ID)
	}
	err := q.Select()
	return users, err
}

// Delete sets deleted_at for a user
func (u User) Delete(db orm.DB, user stems.User) error {
	return db.Delete(&user)
}
