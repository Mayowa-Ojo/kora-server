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
	ID         primitive.ObjectID   `json:"id,omitempty"           bson:"_id,omitempty"`
	Username   string               `json:"username,omitempty"     bson:"username"`
	Firstname  string               `json:"firstname,omitempty"    bson:"firstname"`
	Lastname   string               `json:"lastname,omitempty"     bson:"lastname"`
	Email      string               `json:"email,omitempty"        bson:"email"`
	Hash       string               `json:"hash,omitempty"         bson:"hash"`
	About      string               `json:"about,omitempty"        bson:"about"`
	Credential string               `json:"credential,omitempty"   bson:"credential"`
	Avatar     string               `json:"avatar,omitempty"       bson:"avatar"`
	Views      int                  `json:"views"                  bson:"views"`
	Shares     []primitive.ObjectID `json:"-"                      bson:"shares"`
	Upvotes    int                  `json:"upvotes"                bson:"upvotes"`
	Downvotes  int                  `json:"downvotes"              bson:"downvotes"`
	Followers  []primitive.ObjectID `json:"-"                      bson:"followers"`
	Following  []primitive.ObjectID `json:"-"                      bson:"following"`
	Answers    []primitive.ObjectID `json:"-"                      bson:"answers"`
	Questions  []primitive.ObjectID `json:"-"                      bson:"questions"`
	Posts      []primitive.ObjectID `json:"-"                      bson:"posts"`
	Knowledge  []Topic              `json:"knowledge,omitempty"    bson:"knowledge"`
	Spaces     []primitive.ObjectID `json:"spaces,omitempty"       bson:"spaces"`
	PinnedPost primitive.ObjectID   `json:"pinnedPost,omitempty"   bson:"pinned_post"`
	CreatedAt  time.Time            `json:"createdAt"              bson:"created_at"`
	UpdatedAt  time.Time            `json:"updatedAt"              bson:"updated_at"`
	DeletedAt  time.Time            `json:"deletedAt"              bson:"deleted_at"`
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

// SetDefaultValues - set default values <[]> to array fields instead of <nil>
func (u *User) SetDefaultValues() {
	u.Followers = []primitive.ObjectID{}
	u.Following = []primitive.ObjectID{}
	u.Answers = []primitive.ObjectID{}
	u.Questions = []primitive.ObjectID{}
	u.Posts = []primitive.ObjectID{}
	u.Spaces = []primitive.ObjectID{}
	u.Knowledge = []Topic{}
	u.Avatar = "https://kora-s3-bucket.s3.us-east-2.amazonaws.com/images/default-user-avatar.png"
}
