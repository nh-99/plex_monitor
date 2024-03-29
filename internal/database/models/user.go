package models

import (
	"context"
	"plex_monitor/internal/database"
	"plex_monitor/internal/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	// AnonymousUserID is the ID of the anonymous user.
	AnonymousUserID = "ANONYMOUS"
	// SystemUserID is the ID of the system user.
	SystemUserID = "__SYSTEM__"
)

// User is the struct that represents the user data that is stored in the database.
type User struct {
	ID             string     `bson:"id"`
	Email          string     `bson:"email"`
	Password       []byte     `bson:"-"`
	HashedPassword string     `bson:"password"`
	Activated      bool       `bson:"activated"`
	CreatedAt      time.Time  `bson:"created_at"`
	CreatedBy      string     `bson:"created_by"`
	UpdatedAt      time.Time  `bson:"updated_at"`
	UpdatedBy      string     `bson:"updated_by"`
	DeletedAt      *time.Time `bson:"deleted_at,omitempty"`
	DeletedBy      *string    `bson:"deleted_by,omitempty"`
}

// IsAnonymous checks if the user is the anonymous user.
func (u User) IsAnonymous() bool {
	return u.ID == AnonymousUserID
}

// GetAnonymousUser returns the anonymous user.
func GetAnonymousUser() User {
	return User{ID: AnonymousUserID}
}

// GetUser returns the user with the supplied ID or email.
func GetUser(id string, email string) (User, error) {
	var user User

	var filter bson.M
	if id != "" {
		filter = bson.M{"id": id}
	} else if email != "" {
		filter = bson.M{"email": email}
	}

	err := database.DB.Collection("users").FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// CheckPassword checks that the supplied password and the hashed password in the database match.
func (u User) CheckPassword(password string) bool {
	return utils.CompareStringToHash(password, u.HashedPassword)
}
