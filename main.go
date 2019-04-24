package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jnst/event-driven-architecture/pubsub"
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

	fmt.Println("========== Prepare done ==========")

	// Subscriber polling topic every second.
	// In this case, you only receive messages by polling, but in a real application, processing is performed according to the business domain.
	subscriber := pubsub.NewSubscriber(butler.Sqs)
	ctx := context.Background()
	go subscriber.Subscribe(ctx, broker.QueueUrl)

	// Publisher sends message to topic every 5 seconds.
	publisher := pubsub.NewPublisher(butler.Sns)
	for i := 1; i < 5; i++ {
		_, _ = publisher.Publish(broker.TopicArn, fmt.Sprintf("test-message-%d", i))
		time.Sleep(5 * time.Second)
	}

	fmt.Println("========== Sample code done ==========")

	ctx.Done()
	butler.Destroy()
}
