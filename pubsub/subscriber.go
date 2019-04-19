package pubsub

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
)

type Subscriber struct {
	svc *sns.SNS
}

func NewSubscriber(svc *sns.SNS) *Subscriber {
	return &Subscriber{
		svc: svc,
	}
}

func (s *Subscriber) Subscribe(topicArn string) (string, error) {
	output, err := s.svc.Subscribe(&sns.SubscribeInput{
		TopicArn: aws.String(topicArn),
		Protocol: aws.String("http"),
	})
	if err != nil {
		return "", err
	}
	return *output.SubscriptionArn, nil
}
