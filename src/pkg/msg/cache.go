package msg

import (
	"context"
	"errors"
	"log"
	"sync"
)

type MessageHubCacher struct {
	hub   MessageHub
	mutex sync.RWMutex

	producers map[string]MessageHubProducer
	consumers map[string]MessageHubConsumer
}

func NewMessageHubCacher(hub MessageHub) *MessageHubCacher {
	return &MessageHubCacher{
		hub:       hub,
		producers: make(map[string]MessageHubProducer, 0),
		consumers: make(map[string]MessageHubConsumer, 0),
	}
}

func (m *MessageHubCacher) CreateProducer(name string) (MessageHubProducer, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	producer, exists := m.producers[name]
	if exists {
		return producer, nil
	}

	producer, err := m.hub.CreateProducer(name)
	if err != nil {
		return nil, err
	}

	log.Printf("Storing producer '%s'...", name)
	m.producers[name] = producer

	return producer, nil
}

func (m *MessageHubCacher) CreateConsumer(name string) (MessageHubConsumer, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	consumer, exists := m.consumers[name]
	if exists {
		return consumer, nil
	}

	consumer, err := m.hub.CreateConsumer(name)
	if err != nil {
		return nil, err
	}

	log.Printf("Storing consumer '%s'...", name)
	m.consumers[name] = consumer

	return consumer, nil
}

func (m *MessageHubCacher) Cleanup(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var err error

	for n, producer := range m.producers {
		log.Printf("Performing cleanup for producer '%s'...", n)

		if ierr := producer.Cleanup(ctx); ierr != nil {
			err = errors.Join(err, ierr)
		}
	}

	for n, consumer := range m.consumers {
		log.Printf("Performing cleanup for consumer '%s'...", n)

		if ierr := consumer.Cleanup(ctx); ierr != nil {
			err = errors.Join(err, ierr)
		}
	}

	if ierr := m.hub.Cleanup(ctx); ierr != nil {
		err = errors.Join(err, ierr)
	}

	return err
}
