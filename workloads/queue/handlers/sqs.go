package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cs.utexas.edu/zjia/faas-queue/common"
	"cs.utexas.edu/zjia/faas-queue/utils"

	"cs.utexas.edu/zjia/faas/types"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
)

type sqsInitHandler struct {
	env     types.Environment
	awsSess *session.Session
	sqsSvc  *sqs.SQS
}

type sqsProducerHandler struct {
	env     types.Environment
	awsSess *session.Session
	sqsSvc  *sqs.SQS
}

type sqsConsumerHandler struct {
	env     types.Environment
	awsSess *session.Session
	sqsSvc  *sqs.SQS
}

func NewSqsInitHandler(env types.Environment) types.FuncHandler {
	sess := utils.CreateAWSSessionOrDie()
	return &sqsInitHandler{
		env:     env,
		awsSess: sess,
		sqsSvc:  sqs.New(sess),
	}
}

func NewSqsProducerHandler(env types.Environment) types.FuncHandler {
	sess := utils.CreateAWSSessionOrDie()
	return &sqsProducerHandler{
		env:     env,
		awsSess: sess,
		sqsSvc:  sqs.New(sess),
	}
}

func NewSqsConsumerHandler(env types.Environment) types.FuncHandler {
	sess := utils.CreateAWSSessionOrDie()
	return &sqsConsumerHandler{
		env:     env,
		awsSess: sess,
		sqsSvc:  sqs.New(sess),
	}
}

func (h *sqsInitHandler) Call(ctx context.Context, input []byte) ([]byte, error) {
	parsedInput := &common.QueueInitInput{}
	err := json.Unmarshal(input, parsedInput)
	if err != nil {
		return nil, err
	}
	output, err := initSQS(ctx, h.sqsSvc, parsedInput)
	if err != nil {
		return nil, err
	}
	encodedOutput, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}
	return common.CompressData(encodedOutput), nil
}

func (h *sqsProducerHandler) Call(ctx context.Context, input []byte) ([]byte, error) {
	parsedInput := &common.ProducerFnInput{}
	err := json.Unmarshal(input, parsedInput)
	if err != nil {
		return nil, err
	}
	output, err := producerSQS(ctx, h.sqsSvc, parsedInput)
	if err != nil {
		return nil, err
	}
	encodedOutput, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}
	return common.CompressData(encodedOutput), nil
}

func (h *sqsConsumerHandler) Call(ctx context.Context, input []byte) ([]byte, error) {
	parsedInput := &common.ConsumerFnInput{}
	err := json.Unmarshal(input, parsedInput)
	if err != nil {
		return nil, err
	}
	output, err := consumerSQS(ctx, h.sqsSvc, parsedInput)
	if err != nil {
		return nil, err
	}
	encodedOutput, err := json.Marshal(output)
	if err != nil {
		panic(err)
	}
	return common.CompressData(encodedOutput), nil
}

func initSQS(ctx context.Context, svc *sqs.SQS, input *common.QueueInitInput) (*common.FnOutput, error) {
	for _, queueName := range input.QueueNames {
		if err := utils.CreateSQSQueue(svc, queueName); err != nil {
			return &common.FnOutput{
				Success: false,
				Message: fmt.Sprintf("Failed to create queue %s: %v", queueName, err),
			}, nil
		}
	}
	return &common.FnOutput{Success: true}, nil
}

const kDefaultMessageGroupId = "Default-Message-Group-Id"
const kWarmupCount = 100
const kWarmupIntervalFactor = 4

func producerSQS(ctx context.Context, svc *sqs.SQS, input *common.ProducerFnInput) (*common.FnOutput, error) {
	queueUrl, err := utils.SQSGetQueueUrl(svc, input.QueueName)
	if err != nil {
		return &common.FnOutput{
			Success: false,
			Message: fmt.Sprintf("Failed to get queue URL: %v", err),
		}, nil
	}
	isFifoQueue := utils.SQSIsFifoQueue(input.QueueName)
	duration := time.Duration(input.Duration) * time.Second
	interval := time.Duration(input.IntervalMs) * time.Millisecond
	latencies := make([]int, 0, 128)
	startTime := time.Now()
	count := 0
	for time.Since(startTime) < duration {
		payload := utils.RandomString(input.PayloadSize - utils.TimestampStrLen)
		pushStart := time.Now()
		payload = utils.FormatTime(pushStart) + payload
		input := &sqs.SendMessageInput{
			QueueUrl:    aws.String(queueUrl),
			MessageBody: aws.String(payload),
		}
		if isFifoQueue {
			input.MessageGroupId = aws.String(kDefaultMessageGroupId)
			input.MessageDeduplicationId = aws.String(uuid.NewString())
		}
		// pushStart := time.Now()
		_, err := svc.SendMessage(input)
		elapsed := time.Since(pushStart)
		if err != nil {
			return &common.FnOutput{
				Success:  false,
				Message:  fmt.Sprintf("SQS SendMessage failed: %v", err),
				Duration: time.Since(startTime).Seconds(),
			}, nil
		}
		latencies = append(latencies, int(elapsed.Microseconds()))
		count++
		if count < kWarmupCount {
			time.Sleep(pushStart.Add(interval * kWarmupIntervalFactor).Sub(time.Now()))
		} else {
			time.Sleep(pushStart.Add(interval).Sub(time.Now()))
		}
	}
	return &common.FnOutput{
		Success:   true,
		Duration:  time.Since(startTime).Seconds(),
		Latencies: latencies,
	}, nil
}

const kDeleteMessageBatchSize = 10

func consumerSQS(ctx context.Context, svc *sqs.SQS, input *common.ConsumerFnInput) (*common.FnOutput, error) {
	queueUrl, err := utils.SQSGetQueueUrl(svc, input.QueueName)
	if err != nil {
		return &common.FnOutput{
			Success: false,
			Message: fmt.Sprintf("Failed to get queue URL: %v", err),
		}, nil
	}
	isFifoQueue := utils.SQSIsFifoQueue(input.QueueName)
	duration := time.Duration(input.Duration) * time.Second
	interval := time.Duration(input.IntervalMs) * time.Millisecond
	latencies := make([]int, 0, 128)
	// numMessages := make([]int, 0, 128)
	messageHandles := make([]string, 0, 10)
	startTime := time.Now()
	for time.Since(startTime) < duration {
		input := &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueUrl),
			MaxNumberOfMessages: aws.Int64(int64(input.BatchSize)),
		}
		if isFifoQueue {
			input.ReceiveRequestAttemptId = aws.String(uuid.NewString())
		}
		popStart := time.Now()
		result, err := svc.ReceiveMessage(input)
		// elapsed := time.Since(popStart)
		if err != nil {
			return &common.FnOutput{
				Success:  false,
				Message:  fmt.Sprintf("SQS ReceiveMessage failed: %v", err),
				Duration: time.Since(startTime).Seconds(),
			}, nil
		}
		if len(result.Messages) == 0 {
			// ReceiveMessage times out
			continue
		}
		finishTime := time.Now()
		// latencies = append(latencies, int(elapsed.Microseconds()))
		// numMessages = append(numMessages, len(result.Messages))
		if isFifoQueue {
			messageHandles := make([]string, len(result.Messages))
			for idx, message := range result.Messages {
				messageHandles[idx] = *message.ReceiptHandle
			}
			if err := utils.SQSDeleteMessages(svc, queueUrl, messageHandles); err != nil {
				return &common.FnOutput{
					Success:  false,
					Message:  fmt.Sprintf("SQS DeleteMessage failed: %v", err),
					Duration: time.Since(startTime).Seconds(),
				}, nil
			}
		} else {
			for _, message := range result.Messages {
				messageHandles = append(messageHandles, *message.ReceiptHandle)
				startTime := utils.ParseTime(*message.Body)
				elapsed := finishTime.Sub(startTime)
				latencies = append(latencies, int(elapsed.Microseconds()))
				if len(messageHandles) == kDeleteMessageBatchSize {
					err := utils.SQSDeleteMessages(svc, queueUrl, messageHandles)
					if err != nil {
						return &common.FnOutput{
							Success:  false,
							Message:  fmt.Sprintf("SQS DeleteMessage failed: %v", err),
							Duration: time.Since(startTime).Seconds(),
						}, nil
					}
					messageHandles = make([]string, 0, 10)
				}
			}
		}
		time.Sleep(popStart.Add(interval).Sub(time.Now()))
	}
	elapsed := time.Since(startTime)
	if len(messageHandles) > 0 {
		if err := utils.SQSDeleteMessages(svc, queueUrl, messageHandles); err != nil {
			return &common.FnOutput{
				Success:  false,
				Message:  fmt.Sprintf("SQS DeleteMessage failed: %v", err),
				Duration: time.Since(startTime).Seconds(),
			}, nil
		}
	}
	return &common.FnOutput{
		Success:     true,
		Duration:    elapsed.Seconds(),
		Latencies:   latencies,
		// NumMessages: numMessages,
	}, nil
}
