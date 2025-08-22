package gcloud_pubsub_engine

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/googleapis/gax-go/v2/apierror"
	"github.com/staringfun/millsmess/libs/base"
	"github.com/staringfun/millsmess/libs/types"
	"os"
	"time"
)

type GCloudEngineConfig struct {
	ProjectID    string `json:"projectID" yaml:"projectID" env:"GCP_PROJECT_ID" default:"millsmess-local"`
	EmulatorHost string `json:"emulatorHost" yaml:"emulatorHost" env:"GCP_EMULATOR_HOST"`
	Name         string
	InstanceID   string
}

type GCloudEngine struct {
	Config GCloudEngineConfig
	Client *pubsub.Client
	Ctx    context.Context
}

func NewGCloudEngine(config GCloudEngineConfig) *GCloudEngine {
	return &GCloudEngine{Config: config}
}

func (e *GCloudEngine) Connect(ctx context.Context) error {
	if e.Config.EmulatorHost != "" {
		base.LogErrorIfNotNil(ctx, os.Setenv("PUBSUB_EMULATOR_HOST", e.Config.EmulatorHost), "set emulator host")
	}

	if e.Config.ProjectID != "" {
		base.LogErrorIfNotNil(ctx, os.Setenv("PUBSUB_PROJECT_ID", e.Config.ProjectID), "set project id")
	}

	client, err := pubsub.NewClient(ctx, e.Config.ProjectID)
	if err != nil {
		return err
	}
	e.Client = client
	e.Ctx = ctx
	return nil
}

func (e *GCloudEngine) Publish(topic types.TopicName, message base.PubsubMessage, config base.PublishConfig, ctx context.Context) error {
	traceID := base.GetTraceID(ctx)
	if traceID == "" {
		traceID = base.GenerateTraceID()
	}
	attributes := message.Attributes
	if attributes == nil {
		attributes = map[string]string{}
	}
	base.SetTraceIDAttribute(traceID, attributes)
	base.SetInstanceIDAttribute(base.GetInstanceID(ctx), attributes)
	base.SetServiceNameAttribute(base.ServiceNameKey, attributes)

	r := e.Client.Topic(topic.String()).Publish(ctx, &pubsub.Message{
		Data:       message.Data,
		Attributes: attributes,
	})
	if !config.WaitConfirm {
		return nil
	}
	_, err := r.Get(ctx)
	return err
}

func (e *GCloudEngine) FormatSubscriptionName(topic types.TopicName, config base.SubscriptionConfig) string {
	data, _ := json.Marshal(config)
	if config.IsTopic {
		return fmt.Sprintf("%s.%s.%s", topic.String(), e.Config.Name, data)
	}
	return fmt.Sprintf("%s.%s.%s.%s", topic.String(), e.Config.Name, e.Config.InstanceID, data)
}

func (e *GCloudEngine) Subscribe(topic types.TopicName, config base.SubscriptionConfig, f func(base.PubsubMessage, context.Context) error, ctx context.Context) error {
	name := e.FormatSubscriptionName(topic, config)

	err := e.CreateSubscription(topic, config, ctx)
	if err != nil {
		return err
	}

	return e.Client.Subscription(name).Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		traceID := base.GetTraceIDAttribute(msg.Attributes)
		if traceID == "" {
			traceID = base.GenerateTraceID()
		}
		ctx = base.WithTraceID(traceID, ctx)
		ctx = base.GetLogger(ctx).With().Str("traceID", traceID).WithContext(ctx)
		err := f(base.PubsubMessage{
			Data:       msg.Data,
			Attributes: msg.Attributes,
		}, ctx)
		if base.LogErrorIfNotNil(ctx, err, "process msg error") {
			msg.Nack()
		} else {
			msg.Ack()
		}
	})
}

const AlreadyExistsCode = 6

func (e *GCloudEngine) CreateTopic(topic types.TopicName, config base.TopicConfig, ctx context.Context) error {
	_, err := e.Client.CreateTopicWithConfig(ctx, topic.String(), &pubsub.TopicConfig{})
	if err == nil {
		return nil
	}

	var apiError *apierror.APIError
	if !errors.As(err, &apiError) {
		return err
	}
	if apiError.GRPCStatus().Code() != AlreadyExistsCode {
		return err
	}
	return nil
}

const SubscriptionExpirationDuration = time.Hour * 72
const TopicSubscriptionExpirationDuration = time.Hour * 24

func (e *GCloudEngine) CreateSubscription(topic types.TopicName, config base.SubscriptionConfig, ctx context.Context) error {
	name := e.FormatSubscriptionName(topic, config)

	expirationPolicy := SubscriptionExpirationDuration
	if config.IsTopic {
		expirationPolicy = TopicSubscriptionExpirationDuration
	}
	c := pubsub.SubscriptionConfig{
		Topic:            e.Client.Topic(topic.String()),
		ExpirationPolicy: expirationPolicy,
	}
	c.AckDeadline = config.AckDeadline
	if config.RetryConfig.MinBackoff != 0 || config.RetryConfig.MaxBackoff != 0 {
		c.RetryPolicy = &pubsub.RetryPolicy{
			MinimumBackoff: config.RetryConfig.MinBackoff,
			MaximumBackoff: config.RetryConfig.MaxBackoff,
		}
	}
	if config.DeadLetterConfig.Topic != "" || config.DeadLetterConfig.MaxDeliveryAttempts != 0 {
		c.DeadLetterPolicy = &pubsub.DeadLetterPolicy{
			DeadLetterTopic:     config.DeadLetterConfig.Topic.String(),
			MaxDeliveryAttempts: config.DeadLetterConfig.MaxDeliveryAttempts,
		}
	}

	_, err := e.Client.CreateSubscription(ctx, name, c)
	if err == nil {
		return nil
	}

	var apiError *apierror.APIError
	if !errors.As(err, &apiError) {
		return err
	}
	if apiError.GRPCStatus().Code() != AlreadyExistsCode {
		return err
	}

	return nil
}
