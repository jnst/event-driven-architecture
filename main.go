package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	const EnvKeyAwsProfile = "AWS_PROFILE"
	profile, ok := os.LookupEnv(EnvKeyAwsProfile)
	if !ok {
		panic("missing env: " + EnvKeyAwsProfile)
	}

	creds := credentials.NewSharedCredentials("", profile)
	cfg := aws.NewConfig().WithCredentials(creds).WithRegion("ap-northeast-1")
	sess, err := session.NewSession(cfg)
	if err != nil {
		panic(err)
	}

	svc := sqs.New(sess)
	queueUrl := Prepare(svc)

	output1, err := svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String("Hello, World"),
		QueueUrl:    aws.String(queueUrl),
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("send: %+v\n", output1)

	output2, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{QueueUrl: aws.String(queueUrl)})
	fmt.Printf("receive: %+v\n", output2)

	if output2 != nil {
		for _, message := range output2.Messages {
			_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueUrl),
				ReceiptHandle: message.ReceiptHandle,
			})
			if err != nil {
				panic(err)
			}
			fmt.Println("deleted.")
		}
	}
}

func Prepare(svc *sqs.SQS) string {
	output, err := svc.CreateQueue(&sqs.CreateQueueInput{QueueName: aws.String("es-test")})
	if err != nil {
		panic(err)
	}

	return *output.QueueUrl
}

func Destroy(svc *sqs.SQS, queueUrl string) {
	_, err := svc.DeleteQueue(&sqs.DeleteQueueInput{QueueUrl: aws.String(queueUrl)})
	if err != nil {
		panic(err)
	}
}
