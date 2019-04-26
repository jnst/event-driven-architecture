package pubsub

import (
	"context"
	"encoding/json"
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
	usecase := NewUsecase()

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
			if len(output.Messages) == 0 {
				fmt.Println("received: no message")
			} else {
				fmt.Printf("received: %+v\n", output)
			}

			for _, message := range output.Messages {

				fmt.Println("processing:")
				body := MessageBody{}
				err = json.Unmarshal([]byte(*message.Body), &body)
				if err != nil {
					fmt.Println(err)
					continue
				}
				userEvent := UserEvent{}
				err = json.Unmarshal([]byte(body.Message), &userEvent)
				if err != nil {
					fmt.Println(err)
					continue
				}

				switch userEvent.Status {
				case "user.created":
					usecase.DoAnyProcess(userEvent.Status)
				case "user.updated":
					usecase.DoAnyProcess(userEvent.Status)
				case "user.deleted":
					usecase.DoAnyProcess(userEvent.Status)
				default:
					usecase.DoAnyProcess(userEvent.Status)
				}

				// Messages received from the queue need to be deleted at the end of processing.
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
