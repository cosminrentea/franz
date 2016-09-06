package models

import (
	"errors"
	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"os"
	"strings"
)

var (
	DefaultMessageList *MessageManager
)

type Message struct {
	ID    int64  // Unique identifier
	Title string // Description
	Done  bool   // Is this Message done?
}

// NewMessage creates a new message given a title, that can't be empty.
func NewMessage(title string) (*Message, error) {
	if title == "" {
		return nil, errors.New("empty")
	}
	return &Message{0, title, false}, nil
}

type KafkaConfig struct {
	brokers []string
	topic   string
}

// MessageManager manages a list of messages in memory.
type MessageManager struct {
	KafkaConfig
	messages []*Message
	lastID   int64
}

// NewMessageManager returns an empty MessageManager.
func NewMessageManager(kafkaConfig KafkaConfig) *MessageManager {
	return &MessageManager{
		KafkaConfig: kafkaConfig,
	}
}

// Save saves the given Message in the MessageManager.
func (mm *MessageManager) Save(message *Message) error {
	if message.ID == 0 {
		mm.lastID++
		message.ID = mm.lastID
		mm.messages = append(mm.messages, cloneMessage(message))
		return nil
	}

	for i, t := range mm.messages {
		if t.ID == message.ID {
			mm.messages[i] = cloneMessage(message)
			return nil
		}
	}
	return errors.New("unknown message")
}

// cloneMessage creates and returns a deep copy of the given Message.
func cloneMessage(t *Message) *Message {
	c := *t
	return &c
}

// All returns the list of all the Messages in the MessageManager.
func (mm *MessageManager) All() []*Message {
	return mm.messages
}

// Find returns the Message with the given id in the MessageManager and a boolean
// indicating if the id was found.
func (mm *MessageManager) Find(ID int64) (*Message, bool) {
	for _, t := range mm.messages {
		if t.ID == ID {
			return t, true
		}
	}
	return nil, false
}

// Send a message to a Kafka queue. Can return an error.
func (mm *MessageManager) Send(msg *Message) error {
	kafkaProducer, err := sarama.NewSyncProducer(mm.brokers, nil)
	if err != nil {
		log.Error("error when creating Kafka SyncProducer", err)
		return err
	}
	defer func() {
		if errClose := kafkaProducer.Close(); errClose != nil {
			log.Error("error when closing Kafka SyncProducer", errClose)
		}
	}()
	kafkaMessage := &sarama.ProducerMessage{
		Topic: mm.topic,
		Key:   nil,
		Value: sarama.StringEncoder(msg.Title),
	}
	_, _, errSend := kafkaProducer.SendMessage(kafkaMessage)
	return errSend
}

func init() {
	kafkaConfig := KafkaConfig{
		brokers: strings.Split(os.Getenv("FRANZ_BROKERS"), " "),
		topic:   os.Getenv("FRANZ_TOPIC"),
	}
	log.WithField("config", kafkaConfig).Error("init DefaultMessageList")
	DefaultMessageList = NewMessageManager(kafkaConfig)
}
