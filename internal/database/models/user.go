package models

import (
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
	ID               string            `bson:"id"`
	Email            string            `bson:"email"`
	Password         []byte            `bson:"-"`
	HashedPassword   string            `bson:"password"`
	FrontendServices []FrontendService `bson:"frontendServices"`
	Permissions      []PermissionType  `bson:"permissions"`
	Activated        bool              `bson:"activated"`
	CreatedAt        time.Time         `bson:"created_at"`
	CreatedBy        string            `bson:"created_by"`
	UpdatedAt        time.Time         `bson:"updated_at"`
	UpdatedBy        string            `bson:"updated_by"`
	DeletedAt        *time.Time        `bson:"deleted_at,omitempty"`
	DeletedBy        *string           `bson:"deleted_by,omitempty"`
}

// FrontendServiceType is the type of the frontend service.
type FrontendServiceType string

const (
	// FrontendServiceTypeDiscord is the type of the Discord frontend service.
	FrontendServiceTypeDiscord FrontendServiceType = "discord"
	// FrontendServiceTypeWeb is the type of the web frontend service.
	FrontendServiceTypeWeb FrontendServiceType = "web"
)

// FrontendService stores data about a frontend that the user has access to.
type FrontendService struct {
	UserID string              `bson:"id"`
	Type   FrontendServiceType `bson:"type"`
}

// IsSystem checks if the user is the system user.
func (u User) IsSystem() bool {
	return u.ID == SystemUserID
}

// IsAnonymous checks if the user is the anonymous user.
func (u User) IsAnonymous() bool {
	return u.ID == AnonymousUserID
}

// GetAnonymousUser returns the anonymous user.
func GetAnonymousUser() User {
	return User{ID: AnonymousUserID, Email: "anon"}
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

	err := database.DB.Collection("users").FindOne(database.Ctx, filter).Decode(&user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// GetUserWithFrontendUserID returns the user with the supplied frontend service user ID.
func GetUserWithFrontendUserID(frontendServiceUserID string) (*User, error) {
	var user User

	err := database.DB.Collection("users").FindOne(database.Ctx, bson.M{"frontendServices.id": frontendServiceUserID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CheckPassword checks that the supplied password and the hashed password in the database match.
func (u User) CheckPassword(password string) bool {
	return utils.CompareStringToHash(password, u.HashedPassword)
}

// Save saves the user to the database.
func (u User) Save() error {
	_, err := database.DB.Collection("users").UpdateOne(database.Ctx, bson.M{"id": u.ID}, bson.M{"$set": u})
	if err != nil {
		return err
	}

	return nil
}
