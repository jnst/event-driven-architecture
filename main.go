package main

import (
	"os"

	"github.com/jnst/event-sourcing/pubsub"
)

func main() {
	const EnvKeyAwsProfile = "AWS_PROFILE"
	profile, ok := os.LookupEnv(EnvKeyAwsProfile)
	if !ok {
		panic("missing env: " + EnvKeyAwsProfile)
	}

	butler := pubsub.NewButler(profile)
	butler.Prepare("es-topic", "es-queue")
	butler.Destroy()

	//output1, err := svc.SendMessage(&sqs.SendMessageInput{
	//	MessageBody: aws.String("Hello, World"),
	//	QueueUrl:    aws.String(queueUrl),
	//})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("send: %+v\n", output1)
	//
	//output2, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{QueueUrl: aws.String(queueUrl)})
	//fmt.Printf("receive: %+v\n", output2)
	//
	//if output2 != nil {
	//	for _, message := range output2.Messages {
	//		_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
	//			QueueUrl:      aws.String(queueUrl),
	//			ReceiptHandle: message.ReceiptHandle,
	//		})
	//		if err != nil {
	//			panic(err)
	//		}
	//		fmt.Println("deleted.")
	//	}
	//}
}
