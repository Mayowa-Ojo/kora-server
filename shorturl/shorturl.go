package shorturl

import (
	"fmt"
	"time"

	"github.com/Mayowa-Ojo/kora-server/constants"
	"github.com/Mayowa-Ojo/kora-server/types"
	"github.com/Mayowa-Ojo/kora-server/utils"

	"github.com/Mayowa-Ojo/kora-server/config"
	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	conn *config.DBConn
	env  = config.NewEnvConfig()
)

// URL - database schema for URL collection
type URL struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	LongURL   string             `json:"longUrl"      bson:"long_url"`
	ShortURL  string             `json:"shortUrl"     bson:"short_url"`
	ShortCode string             `json:"shortCode"    bson:"short_code"`
	CreatedAt time.Time          `json:"createdAt"    bson:"created_at"`
	UpdatedAt time.Time          `json:"updatedAt"    bson:"updated_at"`
	DeletedAt time.Time          `json:"deletedAt"    bson:"deleted_at"`
}

// InitShortURLService -
func InitShortURLService(app *fiber.App, c *config.DBConn) {
	conn = c
	baseRouter := app.Group("/api/v1")
	router := baseRouter.Group("/urls")

	router.Get("/:shortCode", getURL(conn))
}

// CreateURL - generates a tiny url link for a post
func CreateURL(ctx *fiber.Ctx, payload types.GenericMap) (*URL, error) {

	col := conn.DB.Collection("urls")

	if _, ok := payload["username"]; !ok {
		return nil, constants.ErrUnprocessableEntity
	}
	if _, ok := payload["slug"]; !ok {
		return nil, constants.ErrUnprocessableEntity
	}
	if _, ok := payload["postType"]; !ok {
		return nil, constants.ErrUnprocessableEntity
	}

	shortCode, err := utils.GenerateID()
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	var longURL string
	if payload["postType"] == "answer" {
		longURL = fmt.Sprintf("%s/%s/%s/%s", env.ClientHostname, payload["slug"], payload["postType"], payload["username"])
	} else {
		longURL = fmt.Sprintf("%s/%s/%s", env.ClientHostname, payload["postType"], payload["slug"])
	}
	shortURL := fmt.Sprintf("%s/bit/%s", env.ClientHostname, shortCode)

	instance := &URL{
		LongURL:   longURL,
		ShortURL:  shortURL,
		ShortCode: shortCode,
	}

	insertResult, err := col.InsertOne(ctx.Fasthttp, instance)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	filter := bson.D{{Key: "_id", Value: insertResult.InsertedID}}
	result := col.FindOne(ctx.Fasthttp, filter)

	url := new(URL)
	if err := result.Decode(&url); err != nil {
		return nil, constants.ErrInternalServer
	}

	return url, nil
}

// long-url -> http://localhost/question/why-is-go-such-a-simple-language
// short-url -> http://localhost/bit/435gghf
func getURL(conn *config.DBConn) fiber.Handler {
	return func(ctx *fiber.Ctx) {
		url := new(URL)
		shortCode := ctx.Params("shortCode")
		col := conn.DB.Collection("urls")

		filter := bson.D{{Key: "short_code", Value: shortCode}}
		result := col.FindOne(ctx.Fasthttp, filter)

		if err := result.Decode(&url); err != nil {
			err := constants.ErrInternalServer
			ctx.Next(err)

			return
		}

		time.Sleep(time.Second * 2)

		resp := utils.NewResponse()
		resp.JSONResponse(ctx, true, fiber.StatusOK, "[INFO]: Resource found", url)
	}
}

func deleteURL(conn *config.DBConn) fiber.Handler {
	return func(ctx *fiber.Ctx) {

	}
}
