package sqs

import (
	"context"
	"encoding/json"

	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SqsConfig struct {
	QueueName           string
	Client              *sqs.Client
	MaxNumberOfMessages int
	VisibilityTimeout   int
	WaitTimeSeconds     int
	local               bool
}

// NewSqsClient func to create a new sqs client.
// Inputs: aws.Config, Config
// Output: Spec
func NewSqsClient(configAws aws.Config, cfg Config) Spec {
	sqsClient := sqs.NewFromConfig(configAws)
	return &SqsConfig{
		QueueName:           cfg.QueueName,
		Client:              sqsClient,
		MaxNumberOfMessages: cfg.MaxNumberOfMessages,
		VisibilityTimeout:   cfg.VisibilityTimeout,
		WaitTimeSeconds:     cfg.WaitTimeSeconds,
	}
}

// getQueueURL private func just to get the sqs URL.
// Inputs: SqsConfig struct
func getQueueURL(s SqsConfig) (*sqs.GetQueueUrlOutput, error) {
	sqsInputName := &sqs.GetQueueUrlInput{
		QueueName: aws.String(s.QueueName),
	}
	resutlsqsURL, err := s.Client.GetQueueUrl(context.TODO(), sqsInputName)
	if err != nil {
		return resutlsqsURL, err
	}
	return resutlsqsURL, err
}

// GetSqsMessages func to get all msg[] like all attributes, msg metadata, body. .
// Inputs: SqsConfig struct
// Output: *sqs.ReceiveMessageOutput
func (s SqsConfig) GetSqsMessages() (*sqs.ReceiveMessageOutput, error) {
	if s.local {
		return getLocalMessages(), nil
	}

	visibilityTimeout := s.VisibilityTimeout
	maxNumberOfMessages := s.MaxNumberOfMessages
	waitTimeSeconds := s.WaitTimeSeconds

	resutlsqsURL, err := getQueueURL(s)
	if err != nil {
		return nil, err
	}
	queueURL := resutlsqsURL.QueueUrl

	params := sqs.ReceiveMessageInput{
		MessageAttributeNames: []string{
			string(types.QueueAttributeNameAll),
		},
		QueueUrl:            queueURL,
		MaxNumberOfMessages: int32(maxNumberOfMessages),
		VisibilityTimeout:   int32(visibilityTimeout),
		WaitTimeSeconds:     int32(waitTimeSeconds),
	}
	msg, err := s.Client.ReceiveMessage(context.TODO(), &params)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if msg.Messages != nil {
		log.Printf("[SQS] Msgs Received")
		return msg, err
	}
	return nil, nil
}

func (s SqsConfig) DeleteMessage(msg *string) {
	if s.local {
		return
	}
	// delete messages
	resutlsqsURL, err := getQueueURL(s)
	if err != nil {
		log.Println("ERROR:", err)
	}
	queueURL := resutlsqsURL.QueueUrl

	params := &sqs.DeleteMessageInput{
		QueueUrl:      queueURL,
		ReceiptHandle: msg,
	}
	_, err = s.Client.DeleteMessage(context.TODO(), params)

	if err != nil {
		log.Println("[SQS] error deleting message")
	}
}

// ReadMessages func to read messages from a queue executing a function for each message.
// Inputs: function to execute. It must receive a types.Message as input and return an error.
func (s SqsConfig) ReadMessages(execute func(msg types.Message) error) {
	for {
		res, err := s.GetSqsMessages()
		if err != nil {
			log.Println("ERROR:", err)
		}
		if res != nil {
			for _, msg := range res.Messages {
				err := execute(msg)
				if err != nil {
					log.Println("ERROR:", err)
				} else {
					s.DeleteMessage(msg.ReceiptHandle)
				}
			}
		}
	}
}

// msg
func (s SqsConfig) SendMessage(msg interface{}) (*sqs.SendMessageOutput, error) {
	msgByte, err := json.Marshal(msg)
	if err != nil {
		log.Println("ERROR:", err)
		return nil, err
	}
	str := string(msgByte)
	strPtr := &str
	// get queue url
	resutlsqsURL, err := getQueueURL(s)
	if err != nil {
		log.Println("ERROR:", err)
		return nil, err

	}
	queueURL := resutlsqsURL.QueueUrl

	// load msg imput
	SendmsgImput := &sqs.SendMessageInput{
		QueueUrl:    queueURL,
		MessageBody: strPtr,
	}
	// send msg to the queue
	rst, err := s.Client.SendMessage(context.Background(), SendmsgImput)
	if err != nil {
		log.Println("ERROR:", err)
		return nil, err
	}
	return rst, nil

}
