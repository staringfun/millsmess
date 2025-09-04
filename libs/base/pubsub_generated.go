// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

// Package base: Code generated; DO NOT EDIT;
package base

import (
	"context"
	"github.com/staringfun/millsmess/libs/types"
)

type PubsubRegistry struct {
	Marshaller                Marshaller
	Engine                    PubsubEngine
	MV1ProfilesUpdateRegistry *TypedSubscribers[types.MV1ProfilesUpdate]
	MV1SessionUpdateRegistry  *TypedSubscribers[types.MV1SessionUpdate]
}

func (r *PubsubRegistry) PublishV1ProfilesUpdate(data types.MV1ProfilesUpdate, attributes map[string]string, config PublishConfig, ctx context.Context) error {
	bytes, err := r.Marshaller.Marshal(data)
	if err != nil {
		return err
	}
	if attributes == nil {
		attributes = map[string]string{}
	}
	SetVersionAttribute("1", attributes)
	return r.Engine.Publish(types.TopicNameProfilesUpdate, PubsubMessage{Data: bytes, Attributes: attributes}, config, ctx)
}
func (r *PubsubRegistry) RegisterV1ProfilesUpdateSubscription(ff func(data types.MV1ProfilesUpdate, attributes map[string]string, ctx context.Context) error, config SubscriptionConfig) {
	r.MV1ProfilesUpdateRegistry.RegisterSubscriber(types.TopicNameProfilesUpdate, config, ff)
}

func (r *PubsubRegistry) PublishV1SessionUpdate(data types.MV1SessionUpdate, attributes map[string]string, config PublishConfig, ctx context.Context) error {
	bytes, err := r.Marshaller.Marshal(data)
	if err != nil {
		return err
	}
	if attributes == nil {
		attributes = map[string]string{}
	}
	SetVersionAttribute("1", attributes)
	return r.Engine.Publish(types.TopicNameSessionUpdate, PubsubMessage{Data: bytes, Attributes: attributes}, config, ctx)
}
func (r *PubsubRegistry) RegisterV1SessionUpdateSubscription(ff func(data types.MV1SessionUpdate, attributes map[string]string, ctx context.Context) error, config SubscriptionConfig) {
	r.MV1SessionUpdateRegistry.RegisterSubscriber(types.TopicNameSessionUpdate, config, ff)
}

func (r *PubsubRegistry) HandleProfilesUpdatesMessage(msg PubsubMessage, config SubscriptionConfig, ctx context.Context) error {
	version := GetVersionAttribute(msg.Attributes)
	switch version {
	case "1":
		var data types.MV1ProfilesUpdate
		err := r.Marshaller.Unmarshal(msg.Data, &data)
		if err != nil {
			GetLogger(ctx).Error().Err(err).Msg("unmarshal error")
			return nil
		}
		return r.MV1ProfilesUpdateRegistry.Run(types.TopicNameProfilesUpdate, config, data, msg.Attributes, ctx)
	}
	if config.IsTopic {
		return NotMatchedVersionError
	}
	return nil
}
func (r *PubsubRegistry) HandleSessionUpdatesMessage(msg PubsubMessage, config SubscriptionConfig, ctx context.Context) error {
	version := GetVersionAttribute(msg.Attributes)
	switch version {
	case "1":
		var data types.MV1SessionUpdate
		err := r.Marshaller.Unmarshal(msg.Data, &data)
		if err != nil {
			GetLogger(ctx).Error().Err(err).Msg("unmarshal error")
			return nil
		}
		return r.MV1SessionUpdateRegistry.Run(types.TopicNameSessionUpdate, config, data, msg.Attributes, ctx)
	}
	if config.IsTopic {
		return NotMatchedVersionError
	}
	return nil
}

func (r *PubsubRegistry) GetSubscribers() map[types.TopicName]map[SubscriptionConfig]func(PubsubMessage, context.Context) error {
	topics := map[types.TopicName]map[SubscriptionConfig]func(PubsubMessage, context.Context) error{}
	for _, subscriber := range r.MV1ProfilesUpdateRegistry.subscribers {
		if _, ok := topics[subscriber.Topic]; !ok {
			topics[subscriber.Topic] = map[SubscriptionConfig]func(PubsubMessage, context.Context) error{}
		}
		topics[subscriber.Topic][subscriber.Config] = func(msg PubsubMessage, ctx context.Context) error {
			return r.HandleProfilesUpdatesMessage(msg, subscriber.Config, ctx)
		}
	}
	for _, subscriber := range r.MV1SessionUpdateRegistry.subscribers {
		if _, ok := topics[subscriber.Topic]; !ok {
			topics[subscriber.Topic] = map[SubscriptionConfig]func(PubsubMessage, context.Context) error{}
		}
		topics[subscriber.Topic][subscriber.Config] = func(msg PubsubMessage, ctx context.Context) error {
			return r.HandleSessionUpdatesMessage(msg, subscriber.Config, ctx)
		}
	}
	return topics
}
