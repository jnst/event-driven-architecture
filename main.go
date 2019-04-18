package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
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

	svc := sns.New(sess)

	topicArn, err := Prepare(svc)
	if err != nil {
		panic(err)
	}

	output, err := svc.Publish(&sns.PublishInput{
		Message:  aws.String("test"),
		TopicArn: aws.String(topicArn),
	})

	fmt.Println(output)

	//err = Destroy(svc, topicArn)
	//if err != nil {
	//	panic(err)
	//}
}

func Prepare(svc *sns.SNS) (string, error) {
	output, err := svc.CreateTopic(&sns.CreateTopicInput{
		Name: aws.String("test"),
		//Attributes: map[string]*string{
		//	"": aws.String(""),
		//},
	})
	if err != nil { // FYI: if already topic exists, it's not occurs error.
		return "", err
	}

	return *output.TopicArn, nil
}

func Destroy(svc *sns.SNS, topicArn string) error {
	_, err := svc.DeleteTopic(&sns.DeleteTopicInput{TopicArn: aws.String(topicArn)})
	return err
}
