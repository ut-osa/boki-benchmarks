package utils

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/google/uuid"
)

func CreateAWSSessionOrDie() *session.Session {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatalf("[FATAL] Failed to create AWS session: %v", err)
	}
	return sess
}

func CreateSQSQueue(svc *sqs.SQS, queueName string) error {
	attributes := map[string]*string{
		"MessageRetentionPeriod":        aws.String("600"), // 10 minutes
		"ReceiveMessageWaitTimeSeconds": aws.String("1"),   // 1 second
		"VisibilityTimeout":             aws.String("10"),  // 10 seconds
	}
	if strings.HasSuffix(queueName, ".fifo") {
		attributes["FifoQueue"] = aws.String("true")
		attributes["ContentBasedDeduplication"] = aws.String("false")
		attributes["DeduplicationScope"] = aws.String("messageGroup")
		attributes["FifoThroughputLimit"] = aws.String("perMessageGroupId")
	}
	_, err := svc.CreateQueue(&sqs.CreateQueueInput{
		QueueName:  aws.String(queueName),
		Attributes: attributes,
	})
	return err
}

func SQSGetQueueUrl(svc *sqs.SQS, queueName string) (string, error) {
	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})
	if err != nil {
		return "", err
	}
	return *result.QueueUrl, nil
}

func SQSIsFifoQueue(queueName string) bool {
	return strings.HasSuffix(queueName, ".fifo")
}

func SQSDeleteMessages(svc *sqs.SQS, queueUrl string, handlers []string) error {
	entries := make([]*sqs.DeleteMessageBatchRequestEntry, len(handlers))
	for i, handler := range handlers {
		entries[i] = &sqs.DeleteMessageBatchRequestEntry{
			Id:            aws.String(uuid.NewString()),
			ReceiptHandle: aws.String(handler),
		}
	}
	result, err := svc.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
		Entries:  entries,
		QueueUrl: aws.String(queueUrl),
	})
	if err != nil {
		return err
	}
	if len(result.Failed) > 0 {
		return fmt.Errorf("DeleteMessageBatch failed: %s", *result.Failed[0].Message)
	} else {
		return nil
	}
}
