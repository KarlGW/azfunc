package azfunc

type binding struct {
	Name       string   `json:"name,omitempty"`
	Type       string   `json:"type,omitempty"`
	Direction  string   `json:"direction,omitempty"`
	AuthLevel  string   `json:"authLevel,omitempty"`
	Methods    []string `json:"methods,omitempty"`
	Route      string   `json:"route,omitempty"`
	Connection string   `json:"connection,omitempty"`
	QueueName  string   `json:"queueName,omitempty"`
	TopicName  string   `json:"topicName,omitempty"`
}
