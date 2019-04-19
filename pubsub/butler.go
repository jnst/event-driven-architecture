package pubsub

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Butler struct {
	Sns    *sns.SNS
	Sqs    *sqs.SQS
	broker Broker
}

func NewButler(profile string) *Butler {
	creds := credentials.NewSharedCredentials("", profile)
	cfg := aws.NewConfig().WithCredentials(creds).WithRegion("ap-northeast-1")
	sess, err := session.NewSession(cfg)
	if err != nil {
		panic(err)
	}

	return &Butler{
		Sns: sns.New(sess),
		Sqs: sqs.New(sess),
	}
}

func (b *Butler) Prepare(topicName, queueName string) Broker {
	// Message Queue
	queueOutput, err := b.Sqs.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Prepare queue:\n %+v\n", queueOutput)

	queueAttrOutput, err := b.Sqs.GetQueueAttributes(&sqs.GetQueueAttributesInput{
		AttributeNames: []*string{aws.String("All")},
		QueueUrl:       queueOutput.QueueUrl,
	})
	if err != nil {
		panic(err)
	}
	queueArn := queueAttrOutput.Attributes["QueueArn"]
	fmt.Printf("Prepare queue-attributes:\n %+v\n", queueAttrOutput)

	// Pub/Sub
	topicOutput, err := b.Sns.CreateTopic(&sns.CreateTopicInput{
		Name: aws.String(topicName),
	})
	if err != nil { // FYI: if already topic exists, it's not occurs error.
		panic(err)
	}
	fmt.Printf("Prepare topic:\n %+v\n", topicOutput)

	subOutput, err := b.Sns.Subscribe(&sns.SubscribeInput{
		TopicArn: topicOutput.TopicArn,
		Protocol: aws.String("sqs"),
		Endpoint: queueArn,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Prepare subscription:\n %+v\n", subOutput)

	b.broker = Broker{
		TopicArn:        *topicOutput.TopicArn,
		QueueUrl:        *queueOutput.QueueUrl,
		SubscriptionArn: *subOutput.SubscriptionArn,
	}

	return b.broker
}

func (b *Butler) Destroy() {
	_, err := b.Sns.DeleteTopic(&sns.DeleteTopicInput{TopicArn: aws.String(b.broker.TopicArn)})
	if err != nil {
		panic(err)
	}
	fmt.Println("Destroy topic done.")

	_, err = b.Sqs.DeleteQueue(&sqs.DeleteQueueInput{QueueUrl: aws.String(b.broker.QueueUrl)})
	if err != nil {
		panic(err)
	}
	fmt.Println("Destroy queue done.")
}
