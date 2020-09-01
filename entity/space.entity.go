package entity

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Space - Defines the schema of a space collection
type Space struct {
	ID           primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string               `json:"name"         bson:"name"`
	About        string               `json:"about"        bson:"about"`
	Details      string               `json:"details"      bson:"details"`
	Icon         string               `json:"icon"         bson:"icon"`
	CoverPhoto   string               `json:"coverPhoto"   bson:"cover_photo"`
	Followers    []primitive.ObjectID `json:"-"            bson:"followers"`
	Admins       []User               `json:"admins"       bson:"admins"`
	Moderators   []User               `json:"moderators"   bson:"moderators"`
	Contributors []User               `json:"contributors" bson:"contributors"`
	Posts        []primitive.ObjectID `json:"-"            bson:"posts"`
	CreatedAt    time.Time            `json:"createdAt"    bson:"created_at"`
	UpdatedAt    time.Time            `json:"updatedAt"    bson:"updated_at"`
	DeletedAt    time.Time            `json:"deletedAt"    bson:"deleted_at"`
}

// Validate - validates struct fields against defined rules
func (s Space) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required, validation.Length(1, 25)),
		validation.Field(&s.About, validation.Required, validation.Length(1, 80)),
	)
}
