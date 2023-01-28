package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func NewSQS(s *session.Session) *sqs.SQS {
	return sqs.New(s)
}
