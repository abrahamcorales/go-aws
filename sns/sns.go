package sns

import (
	_ "context"
	"encoding/json"
	"fmt"
	_ "github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"

	aws2 "github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
)

const (
	StringDataType = "String"
	BinaryDataType = "Binary"
)

var (
	ErrPublishMsg = errors.New("error publishing message")
	ErrMarshal    = errors.New("error marshalling item")
)

type SNS struct {
	client   notificationClient
	topicARN string
	local    bool
}

type Config struct {
	ARN    string `json:"topic_name"`
	Region string `json:"region"`
}

type Options struct {
	MessageAttributes map[string]any
	MessageGroupID    *string
	DeduplicationID   *string
}

type notificationClient interface {
	Publish(input *sns.PublishInput) (*sns.PublishOutput, error)
}

type Publisher interface {
	Publish(i any) error
	PublishWithMsgAttributes(i any, ma map[string]any) error
	PublishWithOptions(i any, ops Options) error
}

func (s *SNS) Publish(i any) error {
	body, err := json.Marshal(i)
	if err != nil {
		return errors.Wrap(ErrMarshal, err.Error())
	}

	if s.local {
		return s.publishLocalMsg(body)
	}

	bodyStr := string(body)

	if _, err := s.client.Publish(&sns.PublishInput{
		Message:  &bodyStr,
		TopicArn: &s.topicARN,
	}); err != nil {
		return errors.Wrap(ErrPublishMsg, err.Error())
	}
	return nil
}

func (s *SNS) PublishWithMsgAttributes(i any, ma map[string]any) error {
	awsMa := make(map[string]*sns.MessageAttributeValue)
	for key, value := range ma {
		switch v := value.(type) {
		case string:
			awsMa[key] = &sns.MessageAttributeValue{
				DataType:    aws2.String(StringDataType),
				StringValue: aws2.String(v),
			}
		case []byte:
			awsMa[key] = &sns.MessageAttributeValue{
				DataType:    aws2.String(BinaryDataType),
				BinaryValue: v,
			}
		default:
			awsMa[key] = &sns.MessageAttributeValue{
				DataType:    aws2.String(BinaryDataType),
				BinaryValue: []byte(fmt.Sprintf("%v", v)),
			}
		}
	}

	body, err := json.Marshal(i)
	if err != nil {
		return errors.Wrap(ErrMarshal, err.Error())
	}

	if s.local {
		return s.publishLocalMsg(body)
	}

	bodyStr := string(body)

	if _, err := s.client.Publish(&sns.PublishInput{
		Message:           &bodyStr,
		TopicArn:          &s.topicARN,
		MessageAttributes: awsMa,
	}); err != nil {
		return errors.Wrap(ErrPublishMsg, err.Error())
	}
	return nil
}

func (s *SNS) PublishWithOptions(i any, ops Options) error {
	awsMa := makeMapAttributes(ops.MessageAttributes)

	body, err := json.Marshal(i)
	if err != nil {
		return errors.Wrap(ErrMarshal, err.Error())
	}

	if s.local {
		return s.publishLocalMsg(body)
	}

	bodyStr := string(body)

	if _, err := s.client.Publish(&sns.PublishInput{
		Message:                &bodyStr,
		TopicArn:               &s.topicARN,
		MessageAttributes:      awsMa,
		MessageGroupId:         ops.MessageGroupID,
		MessageDeduplicationId: ops.DeduplicationID,
	}); err != nil {
		return errors.Wrap(ErrPublishMsg, err.Error())
	}
	return nil
}

func makeMapAttributes(ma map[string]any) map[string]*sns.MessageAttributeValue {
	awsMa := make(map[string]*sns.MessageAttributeValue)
	for key, value := range ma {
		switch v := value.(type) {
		case string:
			awsMa[key] = &sns.MessageAttributeValue{
				DataType:    aws2.String(StringDataType),
				StringValue: aws2.String(v),
			}
		case []byte:
			awsMa[key] = &sns.MessageAttributeValue{
				DataType:    aws2.String(BinaryDataType),
				BinaryValue: v,
			}
		default:
			awsMa[key] = &sns.MessageAttributeValue{
				DataType:    aws2.String(BinaryDataType),
				BinaryValue: []byte(fmt.Sprintf("%v", v)),
			}
		}
	}
	return awsMa
}

func NewSNS(cfg *Config) *SNS {
	return &SNS{
		client: sns.New(session.Must(session.NewSessionWithOptions(
			session.Options{
				Config: aws2.Config{
					Region: aws2.String(cfg.Region),
				},
			},
		))),
		topicARN: cfg.ARN,
		local:    false,
	}
}
