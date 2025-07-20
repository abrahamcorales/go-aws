// Package sqs interface func
package sqs

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Spec interface {
	GetSqsMessages() (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(msg *string)
	ReadMessages(execute func(msg types.Message) error)
	SendMessage(msg interface{}) (*sqs.SendMessageOutput, error)
}

type Config struct {
	QueueName           string `json:"queue_name"`
	MaxNumberOfMessages int    `json:"messages_max_number"`
	VisibilityTimeout   int    `json:"visibility_timout"`
	WaitTimeSeconds     int    `json:"wait_time_second"`
}
