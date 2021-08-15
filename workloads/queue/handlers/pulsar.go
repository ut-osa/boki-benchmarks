package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cs.utexas.edu/zjia/faas-queue/common"
	"cs.utexas.edu/zjia/faas-queue/utils"

	"cs.utexas.edu/zjia/faas/types"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/google/uuid"
)

type pulsarProducerHandler struct {
	env    types.Environment
	client pulsar.Client
}

type pulsarConsumerHandler struct {
	env    types.Environment
	client pulsar.Client
}

func NewPulsarProducerHandler(env types.Environment) types.FuncHandler {
	return &pulsarProducerHandler{
		env:    env,
		client: utils.CreatePulsarClientOrDie(),
	}
}

func NewPulsarConsumerHandler(env types.Environment) types.FuncHandler {
	return &pulsarConsumerHandler{
		env:    env,
		client: utils.CreatePulsarClientOrDie(),
	}
}

func (h *pulsarProducerHandler) Call(ctx context.Context, input []byte) ([]byte, error) {
	parsedInput := &common.ProducerFnInput{}
	err := json.Unmarshal(input, parsedInput)
	if err != nil {
		return nil, err
	}
	output, err := producerPulsar(ctx, h.client, parsedInput)
	if err != nil {
		return nil, err
	}
	encodedOutput, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}
	return common.CompressData(encodedOutput), nil
}

func (h *pulsarConsumerHandler) Call(ctx context.Context, input []byte) ([]byte, error) {
	parsedInput := &common.ConsumerFnInput{}
	err := json.Unmarshal(input, parsedInput)
	if err != nil {
		return nil, err
	}
	output, err := consumerPulsar(ctx, h.client, parsedInput)
	if err != nil {
		return nil, err
	}
	encodedOutput, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}
	return common.CompressData(encodedOutput), nil
}

const kTopicPrefix = "persistent://public/default/"

func producerPulsar(ctx context.Context, client pulsar.Client, input *common.ProducerFnInput) (*common.FnOutput, error) {
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: kTopicPrefix + input.QueueName,
	})
	if err != nil {
		return &common.FnOutput{
			Success: false,
			Message: fmt.Sprintf("Failed to create producer: %v", err),
		}, nil
	}
	defer producer.Close()
	duration := time.Duration(input.Duration) * time.Second
	interval := time.Duration(input.IntervalMs) * time.Millisecond
	latencies := make([]int, 0, 128)
	startTime := time.Now()
	for time.Since(startTime) < duration {
		payload := utils.RandomString(input.PayloadSize - utils.TimestampStrLen)
		pushStart := time.Now()
		payload = utils.FormatTime(pushStart) + payload
		message := &pulsar.ProducerMessage{
			Payload:      []byte(payload),
			Key:          uuid.NewString(),
			DeliverAfter: 100 * time.Millisecond,
		}
		_, err := producer.Send(ctx, message)
		elapsed := time.Since(pushStart)
		if err != nil {
			return &common.FnOutput{
				Success:  false,
				Message:  fmt.Sprintf("Producer send failed: %v", err),
				Duration: time.Since(startTime).Seconds(),
			}, nil
		}
		latencies = append(latencies, int(elapsed.Microseconds()))
		time.Sleep(pushStart.Add(interval).Sub(time.Now()))
	}
	return &common.FnOutput{
		Success:   true,
		Duration:  time.Since(startTime).Seconds(),
		Latencies: latencies,
	}, nil
}

const kDefaultSubscriptionName = "Default"

func consumerPulsar(ctx context.Context, client pulsar.Client, input *common.ConsumerFnInput) (*common.FnOutput, error) {
	topics, err := client.TopicPartitions(kTopicPrefix + input.QueueName)
	if err != nil {
		return &common.FnOutput{
			Success: false,
			Message: fmt.Sprintf("Failed to fetch partitions: %v", err),
		}, nil
	}
	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topics:            topics,
		SubscriptionName:  kDefaultSubscriptionName,
		ReceiverQueueSize: 1,
		Type:              pulsar.Shared,
	})
	if err != nil {
		return &common.FnOutput{
			Success: false,
			Message: fmt.Sprintf("Failed to create consumer: %v", err),
		}, nil
	}
	defer consumer.Close()
	duration := time.Duration(input.Duration) * time.Second
	interval := time.Duration(input.IntervalMs) * time.Millisecond
	latencies := make([]int, 0, 128)
	startTime := time.Now()
	for time.Since(startTime) < duration {
		popStart := time.Now()
		newCtx, _ := context.WithTimeout(ctx, 1*time.Second)
		msg, err := consumer.Receive(newCtx)
		// elapsed := time.Since(popStart)
		if err != nil {
			if err == context.DeadlineExceeded {
				continue
			} else if perr, ok := err.(*pulsar.Error); ok && perr.Result() == pulsar.ResultTimeoutError {
				continue
			} else {
				return &common.FnOutput{
					Success:  false,
					Message:  fmt.Sprintf("Consumer receive failed: %v", err),
					Duration: time.Since(startTime).Seconds(),
				}, nil
			}
		}
		payload := string(msg.Payload())
		delay := time.Since(utils.ParseTime(payload))
		latencies = append(latencies, int(delay.Microseconds()))
		consumer.Ack(msg)
		time.Sleep(popStart.Add(interval).Sub(time.Now()))
	}
	elapsed := time.Since(startTime)
	return &common.FnOutput{
		Success:   true,
		Duration:  elapsed.Seconds(),
		Latencies: latencies,
	}, nil
}
