package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"cs.utexas.edu/zjia/faas-queue/common"
	"cs.utexas.edu/zjia/faas-queue/utils"

	"github.com/montanaflynn/stats"
)

var FLAGS_faas_gateway string
var FLAGS_fn_prefix string
var FLAGS_queue_prefix string
var FLAGS_num_queues int
var FLAGS_queue_shards int
var FLAGS_fifo_queues bool
var FLAGS_num_producer int
var FLAGS_num_consumer int
var FLAGS_duration int
var FLAGS_payload_size int
var FLAGS_producer_interval int
var FLAGS_consumer_interval int
var FLAGS_consumer_bsize int
var FLAGS_consumer_fix_shard bool
var FLAGS_blocking_pop bool
var FLAGS_rand_seed int

func init() {
	flag.StringVar(&FLAGS_faas_gateway, "faas_gateway", "127.0.0.1:8081", "")
	flag.StringVar(&FLAGS_fn_prefix, "fn_prefix", "sqs", "")
	flag.StringVar(&FLAGS_queue_prefix, "queue_prefix", "test", "")
	flag.IntVar(&FLAGS_num_queues, "num_queues", 1, "")
	flag.IntVar(&FLAGS_queue_shards, "queue_shards", 1, "")
	flag.BoolVar(&FLAGS_fifo_queues, "fifo_queues", false, "")
	flag.IntVar(&FLAGS_num_producer, "num_producer", 1, "")
	flag.IntVar(&FLAGS_num_consumer, "num_consumer", 1, "")
	flag.IntVar(&FLAGS_duration, "duration", 10, "")
	flag.IntVar(&FLAGS_payload_size, "payload_size", 64, "")
	flag.IntVar(&FLAGS_producer_interval, "producer_interval", 4, "")
	flag.IntVar(&FLAGS_consumer_interval, "consumer_interval", 4, "")
	flag.IntVar(&FLAGS_consumer_bsize, "consumer_bsize", 1, "")
	flag.BoolVar(&FLAGS_consumer_fix_shard, "consumer_fix_shard", false, "")
	flag.BoolVar(&FLAGS_blocking_pop, "blocking_pop", false, "")
	flag.IntVar(&FLAGS_rand_seed, "rand_seed", 23333, "")

	rand.Seed(int64(FLAGS_rand_seed))
}

func invokeProducer(client *http.Client, queueIndex int, response *common.FnOutput, wg *sync.WaitGroup) {
	defer wg.Done()
	queueName := utils.BuildQueueName(FLAGS_queue_prefix, queueIndex, FLAGS_fifo_queues)
	input := &common.ProducerFnInput{
		QueueName:   queueName,
		QueueShards: FLAGS_queue_shards,
		Duration:    FLAGS_duration,
		PayloadSize: FLAGS_payload_size,
		IntervalMs:  FLAGS_producer_interval,
	}
	url := utils.BuildFunctionUrl(FLAGS_faas_gateway, FLAGS_fn_prefix+"QueueProducer")
	if err := utils.JsonPostRequest(client, url, input, response); err != nil {
		log.Printf("[ERROR] Producer request failed: %v", err)
	} else if !response.Success {
		log.Printf("[ERROR] Producer request failed: %s", response.Message)
	}
}

func invokeConsumer(client *http.Client, queueIndex int, shard int, response *common.FnOutput, wg *sync.WaitGroup) {
	defer wg.Done()
	queueName := utils.BuildQueueName(FLAGS_queue_prefix, queueIndex, FLAGS_fifo_queues)
	input := &common.ConsumerFnInput{
		QueueName:   queueName,
		QueueShards: FLAGS_queue_shards,
		FixedShard:  shard,
		Duration:    FLAGS_duration,
		IntervalMs:  FLAGS_consumer_interval,
		BatchSize:   FLAGS_consumer_bsize,
		BlockingPop: FLAGS_blocking_pop,
	}
	url := utils.BuildFunctionUrl(FLAGS_faas_gateway, FLAGS_fn_prefix+"QueueConsumer")
	if err := utils.JsonPostRequest(client, url, input, response); err != nil {
		log.Printf("[ERROR] Consumer request failed: %v", err)
	} else if !response.Success {
		log.Printf("[ERROR] Consumer request failed: %s", response.Message)
	}
}

func printSummary(title string, results []common.FnOutput) {
	latencies := make([]float64, 0, 128)
	tput := float64(0)
	normedLatencies := make([]float64, 0, 128)
	for _, result := range results {
		if result.Success {
			totalMessages := 0
			for idx, elem := range result.Latencies {
				latency := float64(elem) / 1000.0
				latencies = append(latencies, latency)
				if idx < len(result.NumMessages) {
					num := result.NumMessages[idx]
					normedLatencies = append(normedLatencies, latency/float64(num))
					totalMessages += num
				} else {
					totalMessages++
				}
			}
			tput += float64(totalMessages) / result.Duration
		}
	}
	fmt.Printf("[%s]\n", title)
	fmt.Printf("Throughput: %.1f ops per sec\n", tput)
	if len(latencies) > 0 {
		median, _ := stats.Median(latencies)
		p99, _ := stats.Percentile(latencies, 99.0)
		fmt.Printf("Latency: median = %.3fms, tail (p99) = %.3fms\n", median, p99)
	}
	if len(normedLatencies) > 0 {
		median, _ := stats.Median(normedLatencies)
		p99, _ := stats.Percentile(normedLatencies, 99.0)
		fmt.Printf("Normed latency: median = %.3fms, tail (p99) = %.3fms\n", median, p99)
	}
}

func main() {
	flag.Parse()

	if FLAGS_num_producer%FLAGS_num_queues != 0 {
		log.Fatalf("[FATAL] \"num_producer\" must be divisible by \"num_queues\"")
	}
	if FLAGS_num_consumer%FLAGS_num_queues != 0 {
		log.Fatalf("[FATAL] \"num_consumer\" must be divisible by \"num_queues\"")
	}

	if FLAGS_fn_prefix != "sqs" && FLAGS_fifo_queues {
		log.Fatalf("[FATAL] FIFO queues can only be set for SQS functions")
	}
	if FLAGS_queue_shards > 1 && FLAGS_num_queues != 1 {
		log.Fatalf("[FATAL] Only one queue allows for sharded queue")
	}
	if FLAGS_consumer_fix_shard && FLAGS_queue_shards == 1 {
		log.Fatalf("[FATAL] Fix shard can only be set for sharded queue")
	}
	if FLAGS_consumer_fix_shard && FLAGS_num_consumer%FLAGS_queue_shards != 0 {
		log.Fatalf("[FATAL] When fixing shard, \"num_consumer\" must be divisible by \"queue_shards\"")
	}

	client := &http.Client{
		Transport: &http.Transport{
			MaxConnsPerHost: FLAGS_num_producer + FLAGS_num_consumer,
			MaxIdleConns:    FLAGS_num_producer + FLAGS_num_consumer,
			IdleConnTimeout: 30 * time.Second,
		},
		Timeout: time.Duration(FLAGS_duration*2) * time.Second,
	}

	var wg sync.WaitGroup
	producerResults := make([]common.FnOutput, FLAGS_num_producer)
	consumerResults := make([]common.FnOutput, FLAGS_num_consumer)
	for i := 0; i < FLAGS_num_producer; i++ {
		wg.Add(1)
		go invokeProducer(client, i%FLAGS_num_queues, &producerResults[i], &wg)
	}
	for i := 0; i < FLAGS_num_consumer; i++ {
		wg.Add(1)
		shard := -1
		if FLAGS_consumer_fix_shard {
			shard = i % FLAGS_queue_shards
		}
		go invokeConsumer(client, i%FLAGS_num_queues, shard, &consumerResults[i], &wg)
	}

	wg.Wait()
	printSummary("Producer", producerResults)
	printSummary("Consumer", consumerResults)
}
