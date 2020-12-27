package entity

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SpaceSettings -
type SpaceSettings struct {
	CanAddQuestionsAndAnwers bool   `json:"canAddQuestionsAndAnwers" bson:"Can_add_questions_and_anwers"`
	CanAddPosts              bool   `json:"canAddPosts"              bson:"can_add_posts"`
	CanSubmitQuestions       bool   `json:"canSubmitQuestions"        bson:"can_submit_questions"`
	CanSubmitAnswers         bool   `json:"canSubmitAnswers"          bson:"can_submit_answers"`
	CanSubmitPosts           bool   `json:"canSubmitPosts"           bson:"can_submit_posts"`
	WhoCanComment            string `json:"whoCanComment"            bson:"who_can_comment"`
	ModeratorsCanInvite      bool   `json:"moderatorsCanInvite"      bson:"moderators_can_invite"`
	ContributorsCanInvite    bool   `json:"contributorsCanInvite"    bson:"contributors_can_invite"`
	AnyoneCanRequest         bool   `json:"anyoneCanRequest"         bson:"anyone_can_request"`
	AnyoneCanSubmit          bool   `json:"anyoneCanSubmit"          bson:"anyone_can_submit"`
}

// Space - Defines the schema of a space collection
type Space struct {
	ID           primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string               `json:"name"         bson:"name"`
	About        string               `json:"about"        bson:"about"`
	Slug         string               `json:"slug"         bson:"slug"`
	Details      string               `json:"details"      bson:"details"`
	Icon         string               `json:"icon"         bson:"icon"`
	CoverPhoto   string               `json:"coverPhoto"   bson:"cover_photo"`
	Followers    []primitive.ObjectID `json:"followers"    bson:"followers"`
	Admins       []primitive.ObjectID `json:"admins"       bson:"admins"`
	Moderators   []primitive.ObjectID `json:"moderators"   bson:"moderators"`
	Contributors []primitive.ObjectID `json:"contributors" bson:"contributors"`
	Posts        []primitive.ObjectID `json:"posts"        bson:"posts"`
	Topics       []primitive.ObjectID `json:"topics"       bson:"topics"`
	PinnedPost   primitive.ObjectID   `json:"pinnedPost"   bson:"pinned_post"`
	Settings     SpaceSettings        `json:"settings"     bson:"settings"`
	CreatedAt    time.Time            `json:"createdAt"    bson:"created_at"`
	UpdatedAt    time.Time            `json:"updatedAt"    bson:"updated_at"`
	DeletedAt    time.Time            `json:"deletedAt"    bson:"deleted_at"`
}

// Validate - validates struct fields against defined rules
func (s Space) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required, validation.Length(1, 25)),
		validation.Field(&s.About, validation.Required, validation.Length(1, 80)),
		validation.Field(&s.Topics, validation.Length(0, 3)),
	)
}

// SetDefaultValues - set default values <[]> to array fields instead of <nil>
func (s *Space) SetDefaultValues() {
	s.Followers = []primitive.ObjectID{}
	s.Admins = []primitive.ObjectID{}
	s.Moderators = []primitive.ObjectID{}
	s.Contributors = []primitive.ObjectID{}
	s.Posts = []primitive.ObjectID{}
	s.Topics = []primitive.ObjectID{}
	s.Icon = "https://kora-s3-bucket.s3.us-east-2.amazonaws.com/images/default-space-icon.png"
	s.CoverPhoto = "https://kora-s3-bucket.s3.us-east-2.amazonaws.com/images/default-space-cover.png"
}
