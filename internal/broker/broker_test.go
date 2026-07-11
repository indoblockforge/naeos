package broker

import (
	"sync"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	b := NewRedis()

	if b.Name() != "redis" {
		t.Errorf("expected name 'redis', got %s", b.Name())
	}

	err := b.Connect(&Config{Host: "localhost", Port: 6379})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = b.Ping()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var received *Message
	var wg sync.WaitGroup
	wg.Add(1)

	b.Subscribe("test", func(msg *Message) error {
		received = msg
		wg.Done()
		return nil
	})

	msg := NewMessage("test", []byte("hello"))
	b.Publish("test", msg)

	wg.Wait()

	if received == nil {
		t.Fatal("expected message to be received")
	}

	b.Unsubscribe("test")

	err = b.Close()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRabbitMQ(t *testing.T) {
	b := NewRabbitMQ()

	if b.Name() != "rabbitmq" {
		t.Errorf("expected name 'rabbitmq', got %s", b.Name())
	}

	err := b.Connect(&Config{Host: "localhost", Port: 5672})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = b.Ping()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var received *Message
	var wg sync.WaitGroup
	wg.Add(1)

	b.Subscribe("test", func(msg *Message) error {
		received = msg
		wg.Done()
		return nil
	})

	msg := NewMessage("test", []byte("hello"))
	b.Publish("test", msg)

	wg.Wait()

	if received == nil {
		t.Fatal("expected message to be received")
	}

	err = b.Close()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestKafka(t *testing.T) {
	b := NewKafka()

	if b.Name() != "kafka" {
		t.Errorf("expected name 'kafka', got %s", b.Name())
	}

	err := b.Connect(&Config{Host: "localhost", Port: 9092})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = b.Ping()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var received *Message
	var wg sync.WaitGroup
	wg.Add(1)

	b.Subscribe("test", func(msg *Message) error {
		received = msg
		wg.Done()
		return nil
	})

	msg := NewMessage("test", []byte("hello"))
	b.Publish("test", msg)

	wg.Wait()

	if received == nil {
		t.Fatal("expected message to be received")
	}

	err = b.Close()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNotConnected(t *testing.T) {
	b := NewRedis()

	err := b.Publish("test", &Message{})
	if err == nil {
		t.Error("expected error when not connected")
	}

	err = b.Subscribe("test", func(msg *Message) error { return nil })
	if err == nil {
		t.Error("expected error when not connected")
	}

	err = b.Ping()
	if err == nil {
		t.Error("expected error when not connected")
	}
}

func TestManager(t *testing.T) {
	m := NewManager()

	redis := NewRedis()
	rabbitmq := NewRabbitMQ()
	kafka := NewKafka()

	m.Register("redis", redis)
	m.Register("rabbitmq", rabbitmq)
	m.Register("kafka", kafka)

	got, ok := m.Get("redis")
	if !ok {
		t.Fatal("expected broker to be found")
	}
	if got.Name() != "redis" {
		t.Errorf("expected 'redis', got %s", got.Name())
	}

	names := m.List()
	if len(names) != 3 {
		t.Errorf("expected 3 brokers, got %d", len(names))
	}

	m.Remove("redis")
	_, ok = m.Get("redis")
	if ok {
		t.Error("expected broker to be removed")
	}
}

func TestManagerConnectAll(t *testing.T) {
	m := NewManager()

	redis := NewRedis()
	m.Register("redis", redis)

	configs := map[string]*Config{
		"redis": {Host: "localhost", Port: 6379},
	}

	err := m.ConnectAll(configs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestManagerCloseAll(t *testing.T) {
	m := NewManager()

	redis := NewRedis()
	redis.Connect(&Config{})
	m.Register("redis", redis)

	err := m.CloseAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewMessage(t *testing.T) {
	msg := NewMessage("test", []byte("hello"))
	if msg == nil {
		t.Fatal("expected message to be created")
	}
	if msg.Channel != "test" {
		t.Errorf("expected channel 'test', got %s", msg.Channel)
	}
	if string(msg.Payload) != "hello" {
		t.Errorf("expected payload 'hello', got %s", string(msg.Payload))
	}
}

func TestMessageTimestamp(t *testing.T) {
	before := time.Now()
	msg := NewMessage("test", []byte("hello"))
	after := time.Now()

	if msg.Timestamp.Before(before) || msg.Timestamp.After(after) {
		t.Error("expected timestamp between before and after")
	}
}
