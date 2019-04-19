package pubsub

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

type Publisher struct {
	svc *sns.SNS
}

func NewPublisher(svc *sns.SNS) *Publisher {
	return &Publisher{
		svc: svc,
	}
}

func (p *Publisher) Publish(topicArn, message string) (string, error) {
	output, err := p.svc.Publish(&sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(topicArn),
	})
	if err != nil {
		return "", err
	}
	return *output.MessageId, nil
}
