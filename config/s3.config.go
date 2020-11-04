package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// InitAwsSession - creates a new s3 session
func InitAwsSession(env *EnvConfig) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(env.AwsRegion),
		Credentials: credentials.NewStaticCredentials(env.AwsAccessKeyID, env.AwsSecretAccessKey, ""),
	})

	if err != nil {
		return nil, err
	}

	if _, err := sess.Config.Credentials.Get(); err != nil {
		return nil, err
	}

	return sess, nil
}
