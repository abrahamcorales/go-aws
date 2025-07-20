package sqs

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/gofiber/fiber/v2"
)

var localMessages []types.Message

// NewLocalSqsClient creates a new sqs client for local development
// it receives an SQS config and a fiber app
// it will create a new POST route in the app to simulate the SQS queue. The route will be `/sqs/{queueName}`
func NewLocalSqsClient(cfg Config, app *fiber.App) Spec {

	app.Post("/sqs/"+cfg.QueueName, func(c *fiber.Ctx) error {
		localMessages = append(localMessages, types.Message{
			Body: aws.String(string(c.Body())),
		})
		return c.SendString("Message received ok")
	})
	log.Println("Local SQS endpoint started /sqs/" + cfg.QueueName)
	return &SqsConfig{
		QueueName: cfg.QueueName,
		local:     true,
	}
}

func getLocalMessages() *sqs.ReceiveMessageOutput {
	// return messages from local queue
	messages := make([]types.Message, len(localMessages))
	for i, msg := range localMessages {
		messages[i] = msg
	}
	localMessages = nil
	return &sqs.ReceiveMessageOutput{
		Messages: messages,
	}
}
