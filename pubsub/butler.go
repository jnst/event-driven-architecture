package pubsub

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Butler struct {
	Sns    *sns.SNS
	Sqs    *sqs.SQS
	Broker Broker
}

func (b *Butler) NewButler(sns *sns.SNS, sqs *sqs.SQS) *Butler {
	return &Butler{
		Sns: sns,
		Sqs: sqs,
	}
}

func (b *Butler) Prepare(topicName string) Broker {
	output, err := b.Sns.CreateTopic(&sns.CreateTopicInput{
		Name: aws.String(topicName),
	})
	if err != nil { // FYI: if already topic exists, it's not occurs error.
		panic(err)
	}

	b.Broker = Broker{
		TopicArn: *output.TopicArn,
	}

	return b.Broker
}

func (b *Butler) Destroy() {
	_, err := b.Sns.DeleteTopic(&sns.DeleteTopicInput{TopicArn: aws.String(b.Broker.TopicArn)})
	if err != nil {
		panic(err)
	}
}
