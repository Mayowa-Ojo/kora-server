package entity

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Comment - Defines the schema of a comment collection
type Comment struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Content    string             `json:"content"      bson:"content"`
	Author     *User              `json:"author"       bson:"author"`
	ResponseTo primitive.ObjectID `json:"responseTo"   bson:"response_to"`
	Replies    []Comment          `json:"replies"      bson:"replies"`
	Upvotes    int                `json:"upvotes"      bson:"upvotes"`
	Downvotes  int                `json:"downvotes"    bson:"downvotes"`
}

// Validate - validates struct fields against defined rules
func (c Comment) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Content, validation.Required),
		// validation.Field(&c.Author, validation.Required),
		validation.Field(&c.ResponseTo, validation.Required),
	)
}

// SetDefaultValues - set default values <[]> to array fields instead of <nil>
func (c *Comment) SetDefaultValues() {
	c.Replies = []Comment{}
}
