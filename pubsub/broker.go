package pubsub

type Broker struct {
	TopicARN        string
	QueueURL        string
	SubscriptionARN string
}

type MessageBody struct {
	Message   string `json:"message"`
	Type      string `json:"type"`
	TopicARN  string `json:"topicArn"`
	MessageID string `json:"messageId"`
}

type UserEvent struct {
	UserID string `json:"userId"`
	Status string `json:"status"`
	Time   int64  `json:"time"`
}
