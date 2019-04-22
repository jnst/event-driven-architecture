package pubsub

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Subscriber struct {
	svc *sqs.SQS
}

func NewSubscriber(svc *sqs.SQS) *Subscriber {
	return &Subscriber{
		svc: svc,
	}
}

func (s *Subscriber) Subscribe(ctx context.Context, queueUrl string) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			output, err := s.svc.ReceiveMessage(&sqs.ReceiveMessageInput{QueueUrl: aws.String(queueUrl)})
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("received: %+v\n", output)

			for _, message := range output.Messages {
				_, err = s.svc.DeleteMessage(&sqs.DeleteMessageInput{
					QueueUrl:      aws.String(queueUrl),
					ReceiptHandle: message.ReceiptHandle,
				})
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Printf("deleted: %s\n", *message.ReceiptHandle)
			}
		}
	}
}
