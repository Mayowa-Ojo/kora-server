package entity

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post - defines the schema of a post collection
type Post struct {
	ID          primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string               `json:"title"        bson:"title"`
	Content     string               `json:"content"      bson:"content"`
	Upvotes     int                  `json:"upvotes"      bson:"upvotes"`
	Downvotes   int                  `json:"downvotes"    bson:"downvotes"`
	Views       int                  `json:"views"        bson:"views"`
	Shares      int                  `json:"shares"       bson:"shares"`
	SharedBy    []User               `json:"sharedBy"     bson:"shared_by"`
	Tags        []string             `json:"tags"         bson:"tags"`
	UpvotedBy   []primitive.ObjectID `json:"-"            bson:"upvoted_by"`
	DownvotedBy []primitive.ObjectID `json:"-"            bson:"downvoted_by"`
	PostType    string               `json:"postType"     bson:"post_type"`
	Topics      []Topic              `json:"topics"       bson:"topics"`
	Followers   []primitive.ObjectID `json:"-"            bson:"followers"`
	Author      *User                `json:"author"       bson:"author"`
	Answers     []primitive.ObjectID `json:"-"            bson:"answers"`
	Answered    bool                 `json:"answered"     bson:"answered"`
	ResponseTo  primitive.ObjectID   `json:"-"            bson:"response_to"`
	Comments    []Comment            `json:"comments"     bson:"comments"`
	CreatedAt   time.Time            `json:"createdAt"    bson:"created_at"`
	UpdatedAt   time.Time            `json:"updatedAt"    bson:"updated_at"`
	DeletedAt   time.Time            `json:"deletedAt"    bson:"deleted_at"`
}

// Validate - validates struct fields against defined rules
func (p Post) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Title, validation.Length(1, 250)),
		validation.Field(&p.Content, validation.Required),
		validation.Field(&p.PostType, validation.In([]interface{}{"post", "answer", "question"})),
	)
}
