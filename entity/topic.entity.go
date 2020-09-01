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
	Followers []primitive.ObjectID `json:"-"            bson:"followers"`
	Posts     []primitive.ObjectID `json:"-"            bson:"posts"`
	Spaces    []primitive.ObjectID `json:"-"            bson:"spaces"`
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
