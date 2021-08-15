package common

type QueueInitInput struct {
	QueueNames []string `json:"queueNames"`
}

type ProducerFnInput struct {
	QueueName   string `json:"queueName"`
	QueueShards int    `json:"queueShards"`
	Duration    int    `json:"duration"`
	PayloadSize int    `json:"payloadSize"`
	IntervalMs  int    `json:"interval"`
}

type ConsumerFnInput struct {
	QueueName   string `json:"queueName"`
	QueueShards int    `json:"queueShards"`
	FixedShard  int    `json:"fixedShard"`
	Duration    int    `json:"duration"`
	IntervalMs  int    `json:"interval"`
	BatchSize   int    `json:"batchSize"`
	BlockingPop bool   `json:"blocking"`
}

type FnOutput struct {
	Success     bool    `json:"success"`
	Message     string  `json:"message"`
	Duration    float64 `json:"duration"`
	Latencies   []int   `json:"latencies"`
	NumMessages []int   `json:"numMessages"`
}
