package models

import (
	"context"
	"plex_monitor/internal/config"
	"plex_monitor/internal/database"
	"plex_monitor/internal/utils"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User is the struct that represents the user data that is stored in the database.
type User struct {
	ID               primitive.ObjectID `bson:"_id"`
	Email            string             `bson:"email"`
	Password         []byte             `bson:"-"`
	HashedPassword   string             `bson:"password"`
	FrontendServices []FrontendService  `bson:"frontendServices"`
	Permissions      []PermissionType   `bson:"permissions"`
	Activated        bool               `bson:"activated"`
	CreatedAt        time.Time          `bson:"created_at"`
	CreatedBy        primitive.ObjectID `bson:"created_by"`
	UpdatedAt        time.Time          `bson:"updated_at"`
	UpdatedBy        primitive.ObjectID `bson:"updated_by"`
	DeletedAt        *time.Time         `bson:"deleted_at,omitempty"`
	DeletedBy        *string            `bson:"deleted_by,omitempty"`
}

// FrontendServiceType is the type of the frontend service.
type FrontendServiceType string

// key is used to store/retrieve a User from a context.Context.
type key string

const (
	// FrontendServiceTypeDiscord is the type of the Discord frontend service.
	FrontendServiceTypeDiscord FrontendServiceType = "discord"
	// FrontendServiceTypeWeb is the type of the web frontend service.
	FrontendServiceTypeWeb FrontendServiceType = "web"
	// ContextKeyUserID is the key used to the the user struct in the HTTP context
	ContextKeyUserID key = "user"
)

// FrontendService stores data about a frontend that the user has access to.
type FrontendService struct {
	UserID string              `bson:"id,omitempty"`
	Type   FrontendServiceType `bson:"type,omitempty"`
}

// IsAnonymous checks if the user is the anonymous user.
func (u User) IsAnonymous() bool {
	return u.ID.IsZero()
}

// GetAnonymousUser returns the anonymous user.
func GetAnonymousUser() User {
	logrus.Debug("An anonymous user was returned")
	return User{ID: primitive.NilObjectID, Email: "anon"}
}

// GetUser returns the user with the supplied ID or email.
func GetUser(id string, email string) (User, error) {
	var user User

	var filter bson.M
	if id != "" {
		userID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return GetAnonymousUser(), err
		}
		filter = bson.M{"_id": userID}
	} else if email != "" {
		filter = bson.M{"email": email}
	}

	err := database.DB.Collection("users").FindOne(database.Ctx, filter).Decode(&user)
	if err != nil {
		return GetAnonymousUser(), err
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

// GetBearerToken returns the bearer token for the user.
func (u User) GetBearerToken() (tokenString string, err error) {
	globals := config.GetGlobals()
	_, tokenString, err = globals.JWTAuth.Encode(map[string]interface{}{
		"user_id": u.ID.Hex(),
		"exp":     jwtauth.ExpireIn(1460 * time.Hour), // 1460 hours == two months
	})
	return tokenString, err
}

// Save saves the user to the database.
func (u *User) Save() error {
	if u.ID.IsZero() {
		u.ID = primitive.NewObjectID()
		u.CreatedAt = time.Now()
	}
	opts := options.Update().SetUpsert(true)
	u.UpdatedAt = time.Now()
	_, err := database.DB.Collection("users").UpdateOne(database.Ctx, bson.M{"_id": u.ID}, bson.M{"$set": u}, opts)
	if err != nil {
		return err
	}

	return nil
}

// Reload reloads the user from the database.
func (u *User) Reload() error {
	err := database.DB.Collection("users").FindOne(database.Ctx, bson.M{"_id": u.ID}).Decode(u)
	if err != nil {
		return err
	}
	return nil
}

// NewContext returns a new Context that carries value u.
func (u User) NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ContextKeyUserID, u)
}

// UserFromContext returns the User value stored in ctx, if any.
func UserFromContext(ctx context.Context) (User, bool) {
	u, ok := ctx.Value(ContextKeyUserID).(User)
	return u, ok
}
