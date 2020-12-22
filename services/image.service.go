package services

import (
	"strings"

	"github.com/Mayowa-Ojo/kora/constants"
	"github.com/Mayowa-Ojo/kora/domain"
	"github.com/Mayowa-Ojo/kora/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gofiber/fiber"
)

// var (
// 	env = config.NewEnvConfig() [env declared in auth.service]
// )

// ImageService -
type ImageService struct {
	postRepo  domain.PostRepository
	userRepo  domain.UserRepository
	spaceRepo domain.SpaceRepository
	sess      *session.Session
}

// NewImageService -
func NewImageService(
	p domain.PostRepository,
	u domain.UserRepository,
	s domain.SpaceRepository,
	sess *session.Session,
) domain.ImageService {
	return &ImageService{
		p,
		u,
		s,
		sess,
	}
}

// UploadImage - upload an image file to s3 bucket
//             - returns image url
func (i *ImageService) UploadImage(ctx *fiber.Ctx) (types.GenericMap, error) {
	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		return nil, constants.ErrUnprocessableEntity
	}

	file, err := fileHeader.Open()
	defer file.Close()
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	uploader := s3manager.NewUploader(i.sess)
	slice := strings.Split(fileHeader.Filename, ".")
	ext := slice[len(slice)-1]
	contentType := "image/" + ext
	input := &s3manager.UploadInput{
		Bucket:      aws.String(env.AwsS3Bucket),
		Key:         aws.String("images/" + fileHeader.Filename),
		Body:        file,
		ContentType: &contentType,
	}

	output, err := uploader.Upload(input)
	if err != nil {
		return nil, constants.ErrInternalServer
	}

	result := types.GenericMap{
		"url": output.Location,
	}

	return result, nil
}

// GetImage -
func (i *ImageService) GetImage(ctx *fiber.Ctx) error {
	return nil
}

// DeleteImage -
func (i *ImageService) DeleteImage(ctx *fiber.Ctx) error {
	return nil
}
