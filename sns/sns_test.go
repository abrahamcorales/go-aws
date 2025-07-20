package sns

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/stretchr/testify/assert"
)

type notificationClientMock struct {
	funcPublish func(input *sns.PublishInput) (*sns.PublishOutput, error)
}

func (s *notificationClientMock) Publish(input *sns.PublishInput) (*sns.PublishOutput, error) {
	return s.funcPublish(input)
}

func TestSNS_Publish(t *testing.T) {
	a := assert.New(t)

	t.Run("Success", func(t *testing.T) {
		snsMock := SNS{
			client: &notificationClientMock{
				funcPublish: func(input *sns.PublishInput) (*sns.PublishOutput, error) {
					return nil, nil
				}},
			topicARN: "",
		}

		err := snsMock.Publish("")

		a.Nil(err)
	})

	t.Run("Error marshal", func(t *testing.T) {
		snsMock := SNS{
			client: &notificationClientMock{
				funcPublish: func(input *sns.PublishInput) (*sns.PublishOutput, error) {
					return nil, nil
				}},
			topicARN: "",
		}

		err := snsMock.Publish(make(chan int))

		a.NotNil(err)
		a.ErrorIs(err, ErrMarshal)
	})

	t.Run("Error Publish", func(t *testing.T) {
		snsMock := SNS{
			client: &notificationClientMock{
				funcPublish: func(input *sns.PublishInput) (*sns.PublishOutput, error) {
					return nil, errors.New("error")
				}},
			topicARN: "",
		}

		err := snsMock.Publish("")

		a.NotNil(err)
		a.ErrorIs(err, ErrPublishMsg)
	})
}

func TestSNS_PublishWithMsgAttributes(t *testing.T) {
	a := assert.New(t)

	t.Run("Success string", func(t *testing.T) {
		snsMock := SNS{
			client: &notificationClientMock{
				funcPublish: func(input *sns.PublishInput) (*sns.PublishOutput, error) {
					return nil, nil
				}},
			topicARN: "",
		}

		err := snsMock.PublishWithMsgAttributes("", map[string]any{"test": "test"})
		a.Nil(err)
	})

	t.Run("Success byte array", func(t *testing.T) {
		snsMock := SNS{
			client: &notificationClientMock{
				funcPublish: func(input *sns.PublishInput) (*sns.PublishOutput, error) {
					return nil, nil
				}},
			topicARN: "",
		}

		err := snsMock.PublishWithMsgAttributes("", map[string]any{"test": []byte{}})
		a.Nil(err)
	})

	t.Run("Success struct", func(t *testing.T) {
		snsMock := SNS{
			client: &notificationClientMock{
				funcPublish: func(input *sns.PublishInput) (*sns.PublishOutput, error) {
					return nil, nil
				}},
			topicARN: "",
		}

		err := snsMock.PublishWithMsgAttributes("", map[string]any{"test": struct{}{}})
		a.Nil(err)
	})

	t.Run("Error marshal", func(t *testing.T) {
		snsMock := SNS{
			client: &notificationClientMock{
				funcPublish: func(input *sns.PublishInput) (*sns.PublishOutput, error) {
					return nil, nil
				}},
			topicARN: "",
		}

		err := snsMock.PublishWithMsgAttributes(make(chan int), map[string]any{})

		a.NotNil(err)
		a.ErrorIs(err, ErrMarshal)
	})

	t.Run("Error Publish", func(t *testing.T) {
		snsMock := SNS{
			client: &notificationClientMock{
				funcPublish: func(input *sns.PublishInput) (*sns.PublishOutput, error) {
					return nil, errors.New("error")
				}},
			topicARN: "",
		}

		err := snsMock.PublishWithMsgAttributes("", map[string]any{})

		a.NotNil(err)
		a.ErrorIs(err, ErrPublishMsg)
	})
}

func TestSNS_PublishWithOptions(t *testing.T) {
	a := assert.New(t)

	t.Run("Success string with ops", func(t *testing.T) {
		snsMock := SNS{
			client: &notificationClientMock{
				funcPublish: func(input *sns.PublishInput) (*sns.PublishOutput, error) {
					return nil, nil
				}},
			topicARN: "",
		}

		ops := Options{
			MessageAttributes: map[string]any{"test": "test"},
		}

		err := snsMock.PublishWithOptions("", ops)
		a.Nil(err)
	})

	t.Run("Error marshal with ops", func(t *testing.T) {
		snsMock := SNS{
			client: &notificationClientMock{
				funcPublish: func(input *sns.PublishInput) (*sns.PublishOutput, error) {
					return nil, nil
				}},
			topicARN: "",
		}

		ops := Options{
			MessageAttributes: map[string]any{"test": "test"},
		}

		err := snsMock.PublishWithOptions(make(chan int), ops)

		a.NotNil(err)
		a.ErrorIs(err, ErrMarshal)
	})

	t.Run("Error Publish with ops", func(t *testing.T) {
		snsMock := SNS{
			client: &notificationClientMock{
				funcPublish: func(input *sns.PublishInput) (*sns.PublishOutput, error) {
					return nil, errors.New("error")
				}},
			topicARN: "",
		}

		ops := Options{
			MessageAttributes: map[string]any{"test": "test"},
		}

		err := snsMock.PublishWithOptions("", ops)

		a.NotNil(err)
		a.ErrorIs(err, ErrPublishMsg)
	})
}
