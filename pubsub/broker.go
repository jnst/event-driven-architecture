package pubsub

type Broker struct {
	TopicArn        string
	QueueUrl        string
	SubscriptionArn string
}

type MessageBody struct {
	Message   string `json:"message"`
	Type      string `json:"type"`
	TopicArn  string `json:"topicArn"`
	MessageID string `json:"messageId"`
}

type UserEvent struct {
	UserId string `json:"userId"`
	Status string `json:"status"`
	Time   int64  `json:"time"`
}
