package entity

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Topic - Defines the schema of a topic entity
type Topic struct {
	ID        primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string               `json:"name"         bson:"name"`
	Followers []primitive.ObjectID `json:"followers"    bson:"followers"`
	Posts     []primitive.ObjectID `json:"posts"        bson:"posts"`
	Spaces    []primitive.ObjectID `json:"spaces"       bson:"spaces"`
	CreatedAt time.Time            `json:"createdAt"    bson:"created_at"`
	UpdatedAt time.Time            `json:"updatedAt"    bson:"updated_at"`
	DeletedAt time.Time            `json:"deletedAt"    bson:"deleted_at"`
}

// Validate - validates struct fields against defined rules
func (t Topic) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Name, validation.Required, validation.Length(1, 60)),
	)
}

// SetDefaultValues - set default values <[]> to array fields instead of <nil>
func (t *Topic) SetDefaultValues() {
	t.Followers = []primitive.ObjectID{}
	t.Posts = []primitive.ObjectID{}
	t.Spaces = []primitive.ObjectID{}
}
