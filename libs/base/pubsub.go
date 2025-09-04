package base

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/staringfun/millsmess/libs/types"
	"time"
)

func SetTraceIDAttribute(id string, attributes map[string]string) {
	attributes["traceID"] = id
}

func GetTraceIDAttribute(attributes map[string]string) string {
	return attributes["traceID"]
}

func SetInstanceIDAttribute(id string, attributes map[string]string) {
	attributes["instanceID"] = id
}

func GetInstanceIDAttribute(attributes map[string]string) string {
	return attributes["instanceID"]
}

func SetServiceNameAttribute(name string, attributes map[string]string) {
	attributes["service"] = name
}

func GetServiceNameAttribute(attributes map[string]string) string {
	return attributes["service"]
}

func SetVersionAttribute(version string, attributes map[string]string) {
	attributes["version"] = version
}

func GetVersionAttribute(attributes map[string]string) string {
	return attributes["version"]
}

type PublishConfig struct {
	WaitConfirm bool
}

type PubsubMessage struct {
	Data       []byte
	Attributes map[string]string
}

type TopicConfig struct {
}

type SubscriptionRetryConfig struct {
	MinBackoff time.Duration
	MaxBackoff time.Duration
}

type SubscriptionDeadLetterConfig struct {
	Topic               types.TopicName
	MaxDeliveryAttempts int
}

type SubscriptionConfig struct {
	IsTopic          bool
	AckDeadline      time.Duration
	RetryConfig      SubscriptionRetryConfig
	DeadLetterConfig SubscriptionDeadLetterConfig
}

type PubsubEngine interface {
	Connect(context.Context) error
	Publish(types.TopicName, PubsubMessage, PublishConfig, context.Context) error
	Subscribe(types.TopicName, SubscriptionConfig, func(PubsubMessage, context.Context) error, context.Context) error
	CreateTopic(types.TopicName, TopicConfig, context.Context) error
	CreateSubscription(types.TopicName, SubscriptionConfig, context.Context) error
}

type Marshaller interface {
	Marshal(any) ([]byte, error)
	Unmarshal([]byte, any) error
}

type JSONMarshaller struct{}

func (m *JSONMarshaller) Marshal(a any) ([]byte, error) {
	return json.Marshal(a)
}

func (m *JSONMarshaller) Unmarshal(bytes []byte, a any) error {
	return json.Unmarshal(bytes, a)
}

type Pubsub struct {
	PubsubRegistry
}

type SubscriberItem[T any] struct {
	Topic  types.TopicName
	Config SubscriptionConfig
	Func   func(T, map[string]string, context.Context) error
}

type TypedSubscribers[T any] struct {
	subscribers []SubscriberItem[T]
}

func (s *TypedSubscribers[T]) RegisterSubscriber(topic types.TopicName, config SubscriptionConfig, f func(T, map[string]string, context.Context) error) {
	s.subscribers = append(s.subscribers, SubscriberItem[T]{
		Topic:  topic,
		Config: config,
		Func:   f,
	})
}

func (s *TypedSubscribers[T]) Run(topic types.TopicName, config SubscriptionConfig, data T, attributes map[string]string, ctx context.Context) error {
	var err error
	for _, subscriber := range s.subscribers {
		if subscriber.Topic != topic {
			continue
		}
		if subscriber.Config != config {
			continue
		}
		err = subscriber.Func(data, attributes, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

var NotMatchedVersionError = errors.New("version not matched")

func (p *Pubsub) RunSubscribers(ctx context.Context) error {
	topics := p.GetSubscribers()
	for topic := range topics {
		for config := range topics[topic] {
			err := p.Engine.Subscribe(topic, config, topics[topic][config], ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
