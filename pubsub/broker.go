package pubsub

type Broker struct {
	TopicArn        string
	QueueUrl        string
	SubscriptionArn string
}

type MessageBody struct {
	Message   string `json:"Message"`
	Type      string `json:"Type"`
	TopicArn  string `json:"TopicArn"`
	MessageID string `json:"MessageId"`
}

type UserEvent struct {
	UserId string `json:"UserId"`
	Status string `json:"Status"`
	Time   int64  `json:"Time"`
}
