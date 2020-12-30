package entity

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post - defines the schema of a post collection
type Post struct {
	ID               primitive.ObjectID   `json:"id,omitempty"               bson:"_id,omitempty"`
	Title            string               `json:"title,omitempty"            bson:"title"`
	Content          string               `json:"content,omitempty"          bson:"content"`
	ContentTruncated string               `json:"contentTruncated,omitempty" bson:"content_truncated"`
	ContextLink      string               `json:"contextLink,omitempty"      bson:"context_link"`
	Slug             string               `json:"slug,omitempty"             bson:"slug"`
	PostType         string               `json:"postType,omitempty"         bson:"post_type"`
	Upvotes          int                  `json:"upvotes"                    bson:"upvotes"`
	Downvotes        int                  `json:"downvotes"                  bson:"downvotes"`
	Shares           int                  `json:"shares"                     bson:"shares"`
	Views            int                  `json:"views"                      bson:"views"`
	ShareLink        string               `json:"shareLink"                  bson:"share_link"`
	SharedBy         []primitive.ObjectID `json:"sharedBy"                   bson:"shared_by"`
	UpvotedBy        []primitive.ObjectID `json:"upvotedBy"                  bson:"upvoted_by"`
	DownvotedBy      []primitive.ObjectID `json:"downvotedBy"                bson:"downvoted_by"`
	Topics           []primitive.ObjectID `json:"topics"                     bson:"topics"`
	Followers        []primitive.ObjectID `json:"followers"                  bson:"followers"`
	Author           *User                `json:"author"                     bson:"author"`
	Answers          []primitive.ObjectID `json:"answers"                    bson:"answers"`
	ResponseTo       primitive.ObjectID   `json:"responseTo"                 bson:"response_to"`
	Comments         []Comment            `json:"comments"                   bson:"comments"`
	CreatedAt        time.Time            `json:"createdAt"                  bson:"created_at"`
	UpdatedAt        time.Time            `json:"updatedAt"                  bson:"updated_at"`
	DeletedAt        time.Time            `json:"deletedAt"                  bson:"deleted_at"`
}

// SharedPost -
type SharedPost struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Comment   string             `json:"comment"      bson:"comment"`
	Post      *Post              `json:"post"         bson:"post"`
	Author    *User              `json:"author"       bson:"author"`
	Space     *Space             `json:"space"        bson:"space"`
	CreatedAt time.Time          `json:"createdAt"    bson:"created_at"`
	UpdatedAt time.Time          `json:"updatedAt"    bson:"updated_at"`
	DeletedAt time.Time          `json:"deletedAt"    bson:"deleted_at"`
}

// Validate - validates struct fields against defined rules
func (p Post) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Title, validation.When(p.PostType != "post", validation.Required), validation.Length(1, 250)),
		validation.Field(&p.Content, validation.When(p.PostType != "question", validation.Required)),
		validation.Field(&p.Slug, validation.When(p.PostType != "post", validation.Required)),
		validation.Field(&p.ContentTruncated, validation.When(p.PostType != "question", validation.Required)),
		validation.Field(&p.PostType, validation.In("post", "answer", "question"), validation.Required),
	)
}

// Validate - validates struct fields against defined rules
func (s SharedPost) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Comment, validation.Length(1, 250)),
		validation.Field(&s.Post, validation.Required),
		validation.Field(&s.Author, validation.Required),
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
	p.Comments = []Comment{}
}
