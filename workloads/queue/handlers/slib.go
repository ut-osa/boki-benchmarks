package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cs.utexas.edu/zjia/faas-queue/common"
	"cs.utexas.edu/zjia/faas-queue/utils"

	"cs.utexas.edu/zjia/faas/slib/sync"
	"cs.utexas.edu/zjia/faas/types"
)

type slibProducerHandler struct {
	env types.Environment
}

type slibConsumerHandler struct {
	env types.Environment
}

func NewSlibProducerHandler(env types.Environment) types.FuncHandler {
	return &slibProducerHandler{env: env}
}

func NewSlibConsumerHandler(env types.Environment) types.FuncHandler {
	return &slibConsumerHandler{env: env}
}

type QueueIface interface {
	Push(payload string) error
	Pop() (string /* payload */, error)
	PopBlocking() (string /* payload */, error)
}

func (h *slibProducerHandler) Call(ctx context.Context, input []byte) ([]byte, error) {
	parsedInput := &common.ProducerFnInput{}
	err := json.Unmarshal(input, parsedInput)
	if err != nil {
		return nil, err
	}
	output, err := producerSlib(ctx, h.env, parsedInput)
	if err != nil {
		return nil, err
	}
	encodedOutput, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}
	return common.CompressData(encodedOutput), nil
}

func (h *slibConsumerHandler) Call(ctx context.Context, input []byte) ([]byte, error) {
	parsedInput := &common.ConsumerFnInput{}
	err := json.Unmarshal(input, parsedInput)
	if err != nil {
		return nil, err
	}
	output, err := consumerSlib(ctx, h.env, parsedInput)
	if err != nil {
		return nil, err
	}
	encodedOutput, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}
	return common.CompressData(encodedOutput), nil
}

func createQueue(ctx context.Context, env types.Environment, name string, shards int) (QueueIface, error) {
	if shards == 1 {
		return sync.NewQueue(ctx, env, name)
	} else {
		return sync.NewShardedQueue(ctx, env, name, shards)
	}
}

func producerSlib(ctx context.Context, env types.Environment, input *common.ProducerFnInput) (*common.FnOutput, error) {
	duration := time.Duration(input.Duration) * time.Second
	interval := time.Duration(input.IntervalMs) * time.Millisecond
	q, err := createQueue(ctx, env, input.QueueName, input.QueueShards)
	if err != nil {
		return &common.FnOutput{
			Success: false,
			Message: fmt.Sprintf("NewQueue failed: %v", err),
		}, nil
	}
	latencies := make([]int, 0, 128)
	startTime := time.Now()
	for time.Since(startTime) < duration {
		payload := utils.RandomString(input.PayloadSize - utils.TimestampStrLen)
		pushStart := time.Now()
		payload = utils.FormatTime(pushStart) + payload
		err := q.Push(payload)
		elapsed := time.Since(pushStart)
		if err != nil {
			return &common.FnOutput{
				Success:  false,
				Message:  fmt.Sprintf("QueuePush failed: %v", err),
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

func consumerSlib(ctx context.Context, env types.Environment, input *common.ConsumerFnInput) (*common.FnOutput, error) {
	duration := time.Duration(input.Duration) * time.Second
	interval := time.Duration(input.IntervalMs) * time.Millisecond
	// halfInterval := time.Duration(input.IntervalMs/2) * time.Millisecond
	q, err := createQueue(ctx, env, input.QueueName, input.QueueShards)
	if err != nil {
		return &common.FnOutput{
			Success: false,
			Message: fmt.Sprintf("NewQueue failed: %v", err),
		}, nil
	}
	latencies := make([]int, 0, 128)
	startTime := time.Now()
	for time.Since(startTime) < duration {
		var err error
		var payload string
		popStart := time.Now()
		if input.FixedShard != -1 {
			payload, err = q.(*sync.ShardedQueue).PopFromShard(input.FixedShard)
		} else {
			if input.BlockingPop {
				payload, err = q.PopBlocking()
			} else {
				payload, err = q.Pop()
			}
		}
		// elapsed := time.Since(popStart)
		if err != nil {
			if sync.IsQueueEmptyError(err) {
				time.Sleep(popStart.Add(interval).Sub(time.Now()))
				continue
			} else if sync.IsQueueTimeoutError(err) {
				continue
			} else {
				return &common.FnOutput{
					Success:  false,
					Message:  fmt.Sprintf("QueuePop failed: %v", err),
					Duration: time.Since(startTime).Seconds(),
				}, nil
			}
		}
		delay := time.Since(utils.ParseTime(payload))
		latencies = append(latencies, int(delay.Microseconds()))
		time.Sleep(popStart.Add(interval).Sub(time.Now()))
	}
	return &common.FnOutput{
		Success:   true,
		Duration:  time.Since(startTime).Seconds(),
		Latencies: latencies,
	}, nil
}
