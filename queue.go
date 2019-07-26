package main

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/pubsub"
)

// queue pushes events to pubsub topic
type queue struct {
	client *pubsub.Client
	topic  *pubsub.Topic
}

// newQueue is invoked once per Storable life cycle to configure the store
func newQueue(ctx context.Context, projectID, topicName string) (q *queue, err error) {

	if projectID == "" {
		return nil, errors.New("projectID not set")
	}

	if topicName == "" {
		return nil, errors.New("topicName not set")
	}

	if ctx == nil {
		return nil, errors.New("context not set")
	}

	c, e := pubsub.NewClient(ctx, projectID)
	if e != nil {
		return nil, e
	}

	t := c.Topic(topicName)
	topicExists, err := t.Exists(ctx)
	if err != nil {
		return nil, err
	}

	if !topicExists {
		logger.Printf("Topic %s not found, creating...", topicName)
		t, err = c.CreateTopic(ctx, topicName)
		if err != nil {
			return nil, fmt.Errorf("Unable to create topic: %s - %v", topicName, err)
		}
	}

	o := &queue{
		client: c,
		topic:  t,
	}

	return o, nil
}

// push persist the content
func (q *queue) push(ctx context.Context, data []byte) error {
	msg := &pubsub.Message{Data: data}
	result := q.topic.Publish(ctx, msg)
	_, err := result.Get(ctx)
	return err
}
