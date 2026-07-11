package broker

import (
	"fmt"
	"sync"
	"time"
)

// Message Broker Interface

type Broker interface {
	Name() string
	Connect(config *Config) error
	Close() error
	Publish(channel string, msg *Message) error
	Subscribe(channel string, handler MessageHandler) error
	Unsubscribe(channel string) error
	Ping() error
}

type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
	Timeout  time.Duration
}

type Message struct {
	ID        string
	Channel   string
	Payload   []byte
	Timestamp time.Time
}

type MessageHandler func(msg *Message) error

// Redis Adapter

type Redis struct {
	config    *Config
	connected bool
	channels  map[string]MessageHandler
	mu        sync.RWMutex
}

func NewRedis() *Redis {
	return &Redis{
		channels: make(map[string]MessageHandler),
	}
}

func (r *Redis) Name() string {
	return "redis"
}

func (r *Redis) Connect(config *Config) error {
	r.config = config
	r.connected = true
	return nil
}

func (r *Redis) Close() error {
	r.connected = false
	return nil
}

func (r *Redis) Ping() error {
	if !r.connected {
		return fmt.Errorf("not connected")
	}
	return nil
}

func (r *Redis) Publish(channel string, msg *Message) error {
	if !r.connected {
		return fmt.Errorf("not connected")
	}

	r.mu.RLock()
	handler, ok := r.channels[channel]
	r.mu.RUnlock()

	if ok {
		go handler(msg)
	}
	return nil
}

func (r *Redis) Subscribe(channel string, handler MessageHandler) error {
	if !r.connected {
		return fmt.Errorf("not connected")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.channels[channel] = handler
	return nil
}

func (r *Redis) Unsubscribe(channel string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.channels, channel)
	return nil
}

// RabbitMQ Adapter

type RabbitMQ struct {
	config    *Config
	connected bool
	queues    map[string]MessageHandler
	mu        sync.RWMutex
}

func NewRabbitMQ() *RabbitMQ {
	return &RabbitMQ{
		queues: make(map[string]MessageHandler),
	}
}

func (r *RabbitMQ) Name() string {
	return "rabbitmq"
}

func (r *RabbitMQ) Connect(config *Config) error {
	r.config = config
	r.connected = true
	return nil
}

func (r *RabbitMQ) Close() error {
	r.connected = false
	return nil
}

func (r *RabbitMQ) Ping() error {
	if !r.connected {
		return fmt.Errorf("not connected")
	}
	return nil
}

func (r *RabbitMQ) Publish(channel string, msg *Message) error {
	if !r.connected {
		return fmt.Errorf("not connected")
	}

	r.mu.RLock()
	handler, ok := r.queues[channel]
	r.mu.RUnlock()

	if ok {
		go handler(msg)
	}
	return nil
}

func (r *RabbitMQ) Subscribe(channel string, handler MessageHandler) error {
	if !r.connected {
		return fmt.Errorf("not connected")
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.queues[channel] = handler
	return nil
}

func (r *RabbitMQ) Unsubscribe(channel string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.queues, channel)
	return nil
}

// Kafka Adapter

type Kafka struct {
	config    *Config
	connected bool
	topics    map[string]MessageHandler
	mu        sync.RWMutex
}

func NewKafka() *Kafka {
	return &Kafka{
		topics: make(map[string]MessageHandler),
	}
}

func (k *Kafka) Name() string {
	return "kafka"
}

func (k *Kafka) Connect(config *Config) error {
	k.config = config
	k.connected = true
	return nil
}

func (k *Kafka) Close() error {
	k.connected = false
	return nil
}

func (k *Kafka) Ping() error {
	if !k.connected {
		return fmt.Errorf("not connected")
	}
	return nil
}

func (k *Kafka) Publish(channel string, msg *Message) error {
	if !k.connected {
		return fmt.Errorf("not connected")
	}

	k.mu.RLock()
	handler, ok := k.topics[channel]
	k.mu.RUnlock()

	if ok {
		go handler(msg)
	}
	return nil
}

func (k *Kafka) Subscribe(channel string, handler MessageHandler) error {
	if !k.connected {
		return fmt.Errorf("not connected")
	}

	k.mu.Lock()
	defer k.mu.Unlock()
	k.topics[channel] = handler
	return nil
}

func (k *Kafka) Unsubscribe(channel string) error {
	k.mu.Lock()
	defer k.mu.Unlock()
	delete(k.topics, channel)
	return nil
}

// Broker Manager

type Manager struct {
	brokers map[string]Broker
	mu      sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		brokers: make(map[string]Broker),
	}
}

func (m *Manager) Register(name string, broker Broker) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.brokers[name] = broker
}

func (m *Manager) Get(name string) (Broker, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	broker, ok := m.brokers[name]
	return broker, ok
}

func (m *Manager) Remove(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.brokers, name)
}

func (m *Manager) List() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	names := make([]string, 0, len(m.brokers))
	for name := range m.brokers {
		names = append(names, name)
	}
	return names
}

func (m *Manager) ConnectAll(configs map[string]*Config) error {
	for name, config := range configs {
		broker, ok := m.Get(name)
		if !ok {
			continue
		}
		if err := broker.Connect(config); err != nil {
			return fmt.Errorf("failed to connect to %s: %w", name, err)
		}
	}
	return nil
}

func (m *Manager) CloseAll() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for name, broker := range m.brokers {
		if err := broker.Close(); err != nil {
			return fmt.Errorf("failed to close %s: %w", name, err)
		}
	}
	return nil
}

// Message Builder

func NewMessage(channel string, payload []byte) *Message {
	return &Message{
		ID:        generateID(),
		Channel:   channel,
		Payload:   payload,
		Timestamp: time.Now(),
	}
}

func generateID() string {
	return time.Now().Format("20060102150405.000000000")
}
