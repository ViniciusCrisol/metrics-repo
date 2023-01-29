package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func NewS3(s *session.Session) *s3manager.Uploader {
	return s3manager.NewUploader(s)
}
