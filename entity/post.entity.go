package entity

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post - defines the schema of a post collection
type Post struct {
	ID          primitive.ObjectID   `json:"id,omitempty"           bson:"_id,omitempty"`
	Title       string               `json:"title,omitempty"        bson:"title"`
	Content     string               `json:"content,omitempty"      bson:"content"`
	ContextLink string               `json:"ContextLink,omitempty"  bson:"context_link"`
	Upvotes     int                  `json:"upvotes"                bson:"upvotes"`
	Downvotes   int                  `json:"downvotes"              bson:"downvotes"`
	Shares      int                  `json:"shares"                 bson:"shares"`
	SharedBy    []primitive.ObjectID `json:"-"                      bson:"shared_by"`
	UpvotedBy   []primitive.ObjectID `json:"-"                      bson:"upvoted_by"`
	DownvotedBy []primitive.ObjectID `json:"-"                      bson:"downvoted_by"`
	PostType    string               `json:"postType,omitempty"     bson:"post_type"`
	Topics      []primitive.ObjectID `json:"topics,omitempty"       bson:"topics"`
	Followers   []primitive.ObjectID `json:"-"                      bson:"followers"`
	Author      *User                `json:"author,omitempty"       bson:"author"`
	Answers     []primitive.ObjectID `json:"-"                      bson:"answers"`
	ResponseTo  primitive.ObjectID   `json:"-"                      bson:"response_to"`
	Comments    []primitive.ObjectID `json:"-"                      bson:"comments"`
	CreatedAt   time.Time            `json:"createdAt"              bson:"created_at"`
	UpdatedAt   time.Time            `json:"updatedAt"              bson:"updated_at"`
	DeletedAt   time.Time            `json:"deletedAt"              bson:"deleted_at"`
}

// Validate - validates struct fields against defined rules
func (p Post) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Title, validation.When(p.PostType != "post", validation.Required), validation.Length(1, 250)),
		validation.Field(&p.Content, validation.When(p.PostType != "question", validation.Required)),
		validation.Field(&p.PostType, validation.In("post", "answer", "question"), validation.Required),
	)
}

// SetDefaultValues - set default values <[]> to array fields instead of <nil>
func (p *Post) SetDefaultValues() {
	p.SharedBy = []primitive.ObjectID{}
	p.UpvotedBy = []primitive.ObjectID{}
	p.DownvotedBy = []primitive.ObjectID{}
	p.Topics = []primitive.ObjectID{}
	p.Followers = []primitive.ObjectID{}
	p.Answers = []primitive.ObjectID{}
	p.Comments = []primitive.ObjectID{}
}
