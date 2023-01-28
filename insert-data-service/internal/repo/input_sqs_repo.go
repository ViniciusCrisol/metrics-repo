package repo

import (
	"github.com/ViniciusCrisol/metrics-repo/insert-data-service/pkg/input"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type inputSQSRepo struct {
	url string
	sqs *sqs.SQS
}

func NewInputSQSRepo(url string, sqs *sqs.SQS) *inputSQSRepo {
	return &inputSQSRepo{
		url: url,
		sqs: sqs,
	}
}

func (repo *inputSQSRepo) Get() ([]*input.Input, error) {
	r, err := repo.sqs.ReceiveMessage(
		&sqs.ReceiveMessageInput{
			MaxNumberOfMessages: aws.Int64(10),
			VisibilityTimeout:   aws.Int64(60),
			WaitTimeSeconds:     aws.Int64(0),
			QueueUrl:            &repo.url,
		},
	)
	if err != nil {
		// TODO: Log it!
		return nil, err
	}
	return repo.sqsMsgsToInputs(r), nil
}

func (repo *inputSQSRepo) sqsMsgsToInputs(r *sqs.ReceiveMessageOutput) []*input.Input {
	ipts := []*input.Input{}

	for _, m := range r.Messages {
		ipts = append(
			ipts,
			input.NewInput(*m.ReceiptHandle, *m.Body),
		)
	}
	return ipts
}

func (repo *inputSQSRepo) Delete(i *input.Input) error {
	_, err := repo.sqs.DeleteMessage(
		&sqs.DeleteMessageInput{
			QueueUrl:      &repo.url,
			ReceiptHandle: &i.ID,
		},
	)
	if err != nil {
		// TODO: Log it!
		return err
	}
	return nil
}
