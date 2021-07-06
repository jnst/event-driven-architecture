package pubsub

import (
	"encoding/json"
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

func NewButlerWithProfile(profile string) *Butler {
	creds := credentials.NewSharedCredentials("", profile)
	cfg := aws.NewConfig().
		WithCredentials(creds).
		WithRegion("ap-northeast-1")

	sess, err := session.NewSession(cfg)
	if err != nil {
		panic(err)
	}

	return &Butler{
		Sns: sns.New(sess),
		Sqs: sqs.New(sess),
	}
}

func NewButlerWithLocalstack() *Butler {
	cfg := aws.NewConfig().WithRegion("ap-northeast-1")

	snsSess, err := session.NewSession(cfg.WithEndpoint("http://localhost:4566"))
	if err != nil {
		panic(err)
	}

	sqsSess, err := session.NewSession(cfg.WithEndpoint("http://localhost:4566"))
	if err != nil {
		panic(err)
	}

	return &Butler{
		Sns: sns.New(snsSess),
		Sqs: sqs.New(sqsSess),
	}
}

type PolicyDocument struct {
	Version   string
	Id        string
	Statement []StatementEntry
}

type StatementEntry struct {
	Sid       string
	Effect    string
	Principal string
	Action    string
	Resource  string
	Condition Condition
}

type Condition struct {
	ArnEquals map[string]string
}

func (b *Butler) Prepare(topicName, queueName string) Broker {
	// Message Queue
	queueOutput, err := b.Sqs.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Prepare queue:\n%+v\n", queueOutput)

	queueAttrOutput, err := b.Sqs.GetQueueAttributes(&sqs.GetQueueAttributesInput{
		AttributeNames: []*string{aws.String("All")},
		QueueUrl:       queueOutput.QueueUrl,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Prepare queue-attributes:\n%+v\n", queueAttrOutput)

	// Pub/Sub
	topicOutput, err := b.Sns.CreateTopic(&sns.CreateTopicInput{
		Name: aws.String(topicName),
	})
	if err != nil { // FYI: if already topic exists, it's not occurs error.
		panic(err)
	}

	fmt.Printf("Prepare topic:\n%+v\n", topicOutput)

	subOutput, err := b.Sns.Subscribe(&sns.SubscribeInput{
		Endpoint:              queueAttrOutput.Attributes["QueueArn"],
		Protocol:              aws.String("sqs"),
		ReturnSubscriptionArn: aws.Bool(true),
		TopicArn:              topicOutput.TopicArn,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Prepare subscription:\n%+v\n", subOutput)

	// Policy for SQS
	// See https://docs.aws.amazon.com/ja_jp/sns/latest/dg/sns-sqs-as-subscriber.html#SendMessageToSQS.sqs.permissions
	policy := PolicyDocument{
		Version: "2012-10-17",
		Id:      "MyQueuePolicy",
		Statement: []StatementEntry{
			{
				Sid:       "MySQSPolicy001",
				Effect:    "Allow",
				Principal: "*",
				Action:    "sqs:SendMessage",
				Resource:  *queueAttrOutput.Attributes["QueueArn"],
				Condition: Condition{
					ArnEquals: map[string]string{
						"aws:SourceArn": *topicOutput.TopicArn,
					},
				},
			},
		},
	}

	buf, err := json.Marshal(&policy)
	if err != nil {
		panic(err)
	}

	_, err = b.Sqs.SetQueueAttributes(&sqs.SetQueueAttributesInput{
		Attributes: map[string]*string{"Policy": aws.String(string(buf))},
		QueueUrl:   queueOutput.QueueUrl,
	})
	if err != nil {
		panic(err)
	}

	b.broker = Broker{
		TopicArn:        *topicOutput.TopicArn,
		QueueUrl:        *queueOutput.QueueUrl,
		SubscriptionArn: *subOutput.SubscriptionArn,
	}

	return b.broker
}

func (b *Butler) Destroy() {
	_, err := b.Sns.Unsubscribe(&sns.UnsubscribeInput{SubscriptionArn: aws.String(b.broker.SubscriptionArn)})
	if err != nil {
		panic(err)
	}

	fmt.Println("Destroy subscription done.")

	_, err = b.Sns.DeleteTopic(&sns.DeleteTopicInput{TopicArn: aws.String(b.broker.TopicArn)})
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
