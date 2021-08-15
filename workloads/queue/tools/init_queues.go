package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"cs.utexas.edu/zjia/faas-queue/common"
	"cs.utexas.edu/zjia/faas-queue/utils"
)

var FLAGS_faas_gateway string
var FLAGS_fn_prefix string
var FLAGS_queue_prefix string
var FLAGS_num_queues int
var FLAGS_fifo_queues bool

func init() {
	flag.StringVar(&FLAGS_faas_gateway, "faas_gateway", "127.0.0.1:8081", "")
	flag.StringVar(&FLAGS_fn_prefix, "fn_prefix", "sqs", "")
	flag.StringVar(&FLAGS_queue_prefix, "queue_prefix", "test", "")
	flag.IntVar(&FLAGS_num_queues, "num_queues", 1, "")
	flag.BoolVar(&FLAGS_fifo_queues, "fifo_queues", false, "")
}

func main() {
	flag.Parse()

	input := &common.QueueInitInput{
		QueueNames: make([]string, 0, 16),
	}
	for i := 0; i < FLAGS_num_queues; i++ {
		queueName := utils.BuildQueueName(FLAGS_queue_prefix, i, FLAGS_fifo_queues)
		input.QueueNames = append(input.QueueNames, queueName)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	url := utils.BuildFunctionUrl(FLAGS_faas_gateway, FLAGS_fn_prefix+"InitQueue")
	response := &common.FnOutput{}
	if err := utils.JsonPostRequest(client, url, input, response); err != nil {
		log.Printf("[ERROR] InitQueue request failed: %v", err)
	}

	if !response.Success {
		log.Printf("[ERROR] InitQueue failed: %s", response.Message)
	}
}
