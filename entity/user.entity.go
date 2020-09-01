package entity

import (
	"regexp"
	"time"

	"github.com/Mayowa-Ojo/kora/constants"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User - Defines the schema for a user collection
type User struct {
	ID           primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Username     string               `json:"username"     bson:"username"`
	Firstname    string               `json:"firstname"    bson:"firstname"`
	Lastname     string               `json:"lastname"     bson:"lastname"`
	Email        string               `json:"email"        bson:"email"`
	Hash         string               `json:"hash"         bson:"hash"`
	About        string               `json:"about"        bson:"about"`
	Credential   string               `json:"credential"   bson:"credential"`
	Avatar       string               `json:"avatar"       bson:"avatar"`
	Views        int                  `json:"views"        bson:"views"`
	Shares       []primitive.ObjectID `json:"shares"       bson:"shares"`
	Upvotes      int                  `json:"upvotes"      bson:"upvotes"`
	Downvotes    int                  `json:"downvotes"    bson:"downvotes"`
	Followers    []primitive.ObjectID `json:"followers"    bson:"followers"`
	Following    []primitive.ObjectID `json:"following"    bson:"following"`
	Answers      []primitive.ObjectID `json:"answers"      bson:"answers"`
	Questions    []primitive.ObjectID `json:"questions"    bson:"questions"`
	Posts        []primitive.ObjectID `json:"posts"        bson:"posts"`
	Knowledge    []Topic              `json:"knowledge"    bson:"knowledge"`
	Spaces       []Space              `json:"spaces"       bson:"spaces"`
	PinnedAnswer Post                 `json:"pinnedAnswer" bson:"pinned_answer"`
	CreatedAt    time.Time            `json:"createdAt"    bson:"created_at"`
	UpdatedAt    time.Time            `json:"updatedAt"    bson:"updated_at"`
	DeletedAt    time.Time            `json:"deletedAt"    bson:"deleted_at"`
}

// Validate - validates struct fields against defined rules
func (u User) Validate() error {
	emailRule := []validation.Rule{
		validation.Required,
		validation.Match(regexp.MustCompile(constants.EmailRegex)),
	}
	return validation.ValidateStruct(&u,
		validation.Field(&u.Username, validation.Required),
		validation.Field(&u.Firstname, validation.Required),
		validation.Field(&u.Lastname, validation.Required),
		validation.Field(&u.Email, emailRule...),
		validation.Field(&u.Credential, validation.Length(1, 60)),
	)
}
