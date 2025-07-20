
## SQS go pomelo client

### Install
    
```sh
go get github.com/abraham-corales/go-aws/sqs
```

### Basic Usage

#### Client Initialization
```go
client := sqs.NewSqsClient({
	QueueName: "name-sqs",
    Region: "us-east-1",
    MaxNumberOfMessages: 10,
    WaitTimeSeconds: 20,
    VisibilityTimeout: 30,
    PollingInterval: 5,
    MaxRetries: 3,
})

// This executes the given function for each message received from the queue
go sqsClient.ReadMessages(func(msg types.Message) error {
    log.Info("Message received: ", *msg.Body)
    return nil
})
```
#### Local Development
In local development we can use localstack or the bundled `NewLocalSqsClient` function with a fiber application.
This would create POST endpoints for each queue that can be used to send messages to the queue.

Example:
```go
if config.IsLocalEnvironment(env) {
	sqsClient =  sqs.NewLocalSqsClient(sqs.Config{
        QueueName: "test-queue",
    }, app)
	// Now you can make a request to /sqs/test-queue to send a message to the queue
}

go sqsClient.ReadMessages(func(msg types.Message) error {
    log.Info("Message received: ", *msg.Body)
    return nil
})
```