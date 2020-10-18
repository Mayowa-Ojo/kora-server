package services

import (
	"github.com/Mayowa-Ojo/kora/domain"
	"github.com/Mayowa-Ojo/kora/types"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Post - acts as the business logic layer/service for posts
type Post struct {
	postRepo domain.PostRepository
	userRepo domain.UserRepository
}

// NewPostService - creates a new post service instance
func NewPostService(p domain.PostRepository, u domain.UserRepository) domain.PostService {
	return &Post{
		p,
		u,
	}
}

// GetAll - handles business logic to fetch all posts
func (p Post) GetAll(ctx *fiber.Ctx, opts types.GenericMap) ([]types.GenericMap, error) {
	var result []types.GenericMap
	options := options.Find()

	if _, ok := opts["limit"]; ok {
		options.SetLimit(opts["limit"].(int64))
	}
	// fetch all answers
	filter := bson.D{{Key: "post_type", Value: "answer"}}
	answers, err := p.postRepo.GetMany(ctx, filter, options)

	// fetch all questions
	filter = bson.D{{Key: "post_type", Value: "question"}}
	questions, err := p.postRepo.GetMany(ctx, filter, options)

	if err != nil {
		return nil, err
	}
	// fetch followers count
	for _, q := range questions {
		questionMap := make(types.GenericMap, 0)
		followersCount := len(q.Followers)
		answersCount := len(q.Answers)

		questionMap["post"] = q
		questionMap["followersCount"] = followersCount
		questionMap["answersCount"] = answersCount

		result = append(result, questionMap)
	}

	for _, a := range answers {
		answerMap := make(types.GenericMap, 0)
		answerMap["post"] = a

		result = append(result, answerMap)
	}

	// Sort results by date
	return result, nil
}

// GetAllQuestions -
func (p Post) GetAllQuestions(ctx *fiber.Ctx, limit int) {
	return
}
