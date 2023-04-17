package models

import (
	"database/sql"
	database "plex_monitor/internal/database"
	"plex_monitor/internal/utils"
	"time"

	"github.com/google/uuid"
	"gopkg.in/validator.v2"
)

const (
	anonymousUserID = "ANONYMOUS"
)

type User struct {
	ID             string         `json:"id"`
	Email          string         `json:"email" validate:"^[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$"`
	Password       []byte         `json:"-"`
	HashedPassword []byte         `json:"-"`
	Activated      bool           `json:"activated"`
	CreatedAt      time.Time      `json:"created_at"`
	CreatedBy      sql.NullString `json:"created_by"`
	UpdatedAt      time.Time      `json:"updated_at"`
	UpdatedBy      sql.NullString `json:"updated_by"`
	DeletedAt      utils.NullTime `json:"deleted_at"`
	DeletedBy      sql.NullString `json:"deleted_by"`
}

func GetUser(email string, id string) (User, error) {
	// Get a user from their email
	var user = User{}
	SQL := ``
	var querySubstitution string
	if email != "" {
		SQL = `SELECT bin_to_uuid(id) as id, email, password, activated, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by FROM users WHERE email = ? AND activated = true AND deleted_at IS NOT NULL;`
		querySubstitution = email
	} else {
		SQL = `SELECT bin_to_uuid(id) as id, email, password, activated, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by FROM users WHERE id = ? AND activated = true AND deleted_at IS NOT NULL;`
		querySubstitution = id
	}

	err := database.DB.QueryRow(SQL, querySubstitution).Scan(
		&user.ID,
		&user.Email,
		&user.HashedPassword,
		&user.Activated,
		&user.CreatedAt,
		&user.CreatedBy,
		&user.UpdatedAt,
		&user.UpdatedBy,
		&user.DeletedAt,
		&user.DeletedBy,
	)

	if err != nil {
		return user, err
	}
	return user, nil
}

func HardDeleteUser(email string, areYouSure bool) (bool, User, error) {
	// This method is primarily for cleaning up unit tests. The deleted boolean
	// in the database should be used to delete users rather than actually removing
	// them from the database.DB.
	var user = User{}

	if !areYouSure {
		return false, user, nil
	}

	SQL := `DELETE FROM users WHERE email = ? RETURNING bin_to_uuid(id) as id, email, password, activated, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by;`

	err2 := database.DB.QueryRow(SQL, user.Email).Scan(
		&user.ID,
		&user.Email,
		&user.HashedPassword,
		&user.Activated,
		&user.CreatedAt,
		&user.CreatedBy,
		&user.UpdatedAt,
		&user.UpdatedBy,
		&user.DeletedAt,
		&user.DeletedBy,
	)
	if err2 != nil {
		return false, user, err2
	}

	return true, user, nil
}

func ListUsers() ([]*User, error) {
	SQL := `SELECT bin_to_uuid(id) as id, email, password, activated, created_at, created_by FROM users;`
	rows, err := database.DB.Query(SQL)
	if err != nil {
		return nil, err
	}

	users := make([]*User, 0)
	for rows.Next() {
		user := new(User)
		err := rows.Scan(&user.ID, &user.Email, &user.HashedPassword, &user.CreatedAt, &user.CreatedBy)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func NewUser(email string, pw string) (User, error) {
	newUserId, _ := uuid.NewRandom()
	u := User{
		ID:             newUserId.String(),
		Email:          email,
		HashedPassword: hashedPassword(pw),
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	if err := validator.Validate(u); err != nil {
		return u, err
	}

	return u, nil
}

// GetAnonymousUser creates a blank user for anyone that isn't logged in. This allows us to easily check permissions
// through struct methods.
func GetAnonymousUser() User {
	u := User{
		ID:        anonymousUserID,
		Email:     "anonymoose@user.co",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	return u
}

func (u User) Commit() error {
	var err error

	SQL := `INSERT INTO users(id, email, password, created_at) VALUES(uuid_to_bin($1), $2, $3, $4, $5);`

	_, err = database.DB.Exec(SQL, u.ID, u.Email, u.HashedPassword, u.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

// ActivateAccount will activate a user's account so they can login. Typically used for email verification.
func (u User) ActivateAccount() error {
	var err error

	SQL := `UPDATE users SET activated = true WHERE bin_to_uuid(id) = $1;`

	_, err = database.DB.Exec(SQL, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// CheckPassword checks that the supplied password and the hashed password in the database match.
func (u User) CheckPassword(password string) bool {
	hashedPasswordString := utils.BytesToString(u.HashedPassword)
	return utils.CompareStringToHash(password, hashedPasswordString)
}

func hashedPassword(password string) []byte {
	hashedpass, _ := utils.HashString(password)

	return hashedpass
}
