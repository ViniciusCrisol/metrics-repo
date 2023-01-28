package aws

import (
	"github.com/ViniciusCrisol/metrics-repo/insert-data-service/internal/config"
	"github.com/ViniciusCrisol/metrics-repo/insert-data-service/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func NewSession() (*session.Session, error) {
	s, err := session.NewSession(
		&aws.Config{
			Credentials: credentials.NewStaticCredentials(
				config.AWSLogin,
				config.AWSSecret,
				"",
			),
			Region: &config.AWSRegion,
		},
	)
	if err != nil {
		log.Logger.Error(
			"Failed to init AWS session",
			log.Error(err),
		)
		return nil, err
	}
	return s, nil
}
