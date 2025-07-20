# go-aws


**go-aws** is a Go library that simplifies integration with AWS services like DynamoDB, SQS, and SNS, designed to accelerate the development of microservices and cloud-native applications.

> The DynamoDB integration is **multi-table** and **multi-index**: you can work with multiple tables and global/secondary indexes easily from a single client.

## Why use go-aws?
- **Easy to use:** Simple abstractions for the most common AWS services.
- **Ready for local development:** Built-in support for LocalStack and local/mock clients.
- **Extensible:** Designed so you can adapt it to your needs.
- **Open source and maintained:** Perfect for teams looking for speed and best practices.

## Installation

```sh
go get github.com/abraham-corales/go-aws
```

## Supported services
- **DynamoDB:** Read, write, query, multi-table and multi-index support, and local/mock development.
- **SQS:** Send and receive messages, Fiber integration for local endpoints.
- **SNS:** Publish messages, local environment support.

## Quick example

### DynamoDB
```go
import "github.com/abraham-corales/go-aws/dynamodbv2"

// Initialize DynamoDB client with multiple tables
client := dynamodbv2.NewDynamoClientv2(awsConfig,
    dynamodbv2.WithTable(dynamodbv2.DynamoTable{
        TableName: "users",
        PartitionKeyField: "id",
    }),
    dynamodbv2.WithTable(dynamodbv2.DynamoTable{
        TableName: "orders",
        PartitionKeyField: "order_id",
        SortKeyField: "created_at",
        GlobalIndex: "status-index",
    }),
)

// Save an item in the "users" table
err := client.Save("users", user)

// Query using a global index
err = client.QueryGSI("orders", "status-index", query, 10, 1, &results)
```

### SQS
```go
import "github.com/abraham-corales/go-aws/sqs"

sqsClient := sqs.NewSqsClient(awsConfig, sqs.Config{
    QueueName: "my-queue",
})

// Read messages
sqsClient.ReadMessages(func(msg types.Message) error {
    log.Println("Message received:", *msg.Body)
    return nil
})
```

### SNS
```go
import "github.com/abraham-corales/go-aws/sns"

snsClient := sns.NewSNS(&sns.Config{
    ARN: "arn:aws:sns:us-east-1:123456789012:my-topic",
    Region: "us-east-1",
})

err := snsClient.Publish(map[string]string{"hello": "world"})
```

## Local development
- Support for LocalStack and local clients for testing without real AWS.
- Ready-to-use mocks for your tests.

## Contributing
Pull requests and suggestions are welcome! If you find a bug or want to add a feature, open an issue or PR.

---

> Made with ❤️ for the Go community. If you found it useful, leave a star and share it!
