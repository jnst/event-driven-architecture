package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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
	// When the subscriber polling and receives a message, it takes action according to the type of the event.
	subscriber := pubsub.NewSubscriber(butler.Sqs)
	ctx := context.Background()
	go subscriber.Subscribe(ctx, broker.QueueUrl)

	// Publisher sends message to topic every 5 seconds.
	publisher := pubsub.NewPublisher(butler.Sns)
	for i := 1; i < 5; i++ {
		userId := strconv.Itoa(i)
		msg := pubsub.UserEvent{
			UserId: userId,
			Status: "user.created",
			Time:   time.Now().Unix(),
		}
		b, err := json.Marshal(msg)
		if err != nil {
			panic(err)
		}

		_, _ = publisher.Publish(broker.TopicArn, string(b))
		time.Sleep(5 * time.Second)
	}

	fmt.Println("========== Sample code done ==========")

	ctx.Done()
	butler.Destroy()
}
