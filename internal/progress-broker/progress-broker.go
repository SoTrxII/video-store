package progress_broker

import (
	"context"
	"encoding/json"
	"github.com/dapr/go-sdk/client"
)

// ObjectStorage any S3-like storage solution
type ProgressBroker[T PubSubProxy] struct {
	// Name of the Dapr Component to use
	componentName string
	// Name of the topic to publish into
	topic string
	// Client to publish event into
	client *T
	// Current running context
	ctx *context.Context
}

type UploadState int8

const (
	InProgress UploadState = iota
	Done
	Error
)

type UploadInfos struct {
	// Upload job identifier
	JobId string `json:"jobId"`
	// Current state of the upload
	State UploadState `json:"state"`
	Data  interface{} `json:"data"`
}

type PubSubProxy interface {
	PublishEvent(ctx context.Context, pubsubName string, topicName string, data interface{}, opts ...client.PublishEventOption) error
}

type NewBrokerOptions struct {
	Component string
	Topic     string
}

func NewProgressBroker[T PubSubProxy](ctx *context.Context, client *T, opt NewBrokerOptions) (*ProgressBroker[T], error) {
	return &ProgressBroker[T]{
		componentName: opt.Component,
		topic:         opt.Topic,
		client:        client,
		ctx:           ctx,
	}, nil
}

func (eb *ProgressBroker[T]) SendProgress(data UploadInfos) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = (*eb.client).PublishEvent(*eb.ctx, eb.componentName, eb.topic, string(b))
	if err != nil {
		return err
	}
	return nil
}
