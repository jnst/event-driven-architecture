package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jnst/event-sourcing/pubsub"
)

func main() {
	//const EnvKeyAwsProfile = "AWS_PROFILE"
	//profile, ok := os.LookupEnv(EnvKeyAwsProfile)
	//if !ok {
	//	panic("missing env: " + EnvKeyAwsProfile)
	//}
	//butler := pubsub.NewButlerWithProfile(profile)

	butler := pubsub.NewButlerWithLocalstack()
	broker := butler.Prepare("es-topic", "es-queue")

	subscriber := pubsub.NewSubscriber(butler.Sqs)
	ctx := context.Background()
	go subscriber.Subscribe(ctx, broker.QueueUrl)

	publisher := pubsub.NewPublisher(butler.Sns)
	for i := 1; i < 10; i++ {
		_, _ = publisher.Publish(broker.TopicArn, fmt.Sprintf("test-message-%d", i))
		time.Sleep(5 * time.Second)
	}

	ctx.Done()
	butler.Destroy()
}
