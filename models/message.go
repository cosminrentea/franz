package models

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego"
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
		return nil, fmt.Errorf("empty")
	}
	return &Message{0, title, false}, nil
}

// MessageManager manages a list of messages in memory.
type MessageManager struct {
	messages []*Message
	lastID   int64
}

// NewMessageManager returns an empty MessageManager.
func NewMessageManager() *MessageManager {
	return &MessageManager{}
}

// Save saves the given Message in the MessageManager.
func (m *MessageManager) Save(message *Message) error {
	if message.ID == 0 {
		m.lastID++
		message.ID = m.lastID
		m.messages = append(m.messages, cloneMessage(message))
		return nil
	}

	for i, t := range m.messages {
		if t.ID == message.ID {
			m.messages[i] = cloneMessage(message)
			return nil
		}
	}
	return fmt.Errorf("unknown message")
}

// cloneMessage creates and returns a deep copy of the given Message.
func cloneMessage(t *Message) *Message {
	c := *t
	return &c
}

// All returns the list of all the Messages in the MessageManager.
func (m *MessageManager) All() []*Message {
	return m.messages
}

// Find returns the Message with the given id in the MessageManager and a boolean
// indicating if the id was found.
func (m *MessageManager) Find(ID int64) (*Message, bool) {
	for _, t := range m.messages {
		if t.ID == ID {
			return t, true
		}
	}
	return nil, false
}

// Send a message to a Kafka queue. Can return an error.
func (m *MessageManager) Send(message *Message) error {
	kafkaProducer, err := sarama.NewSyncProducer(beego.AppConfig.Strings("KafkaBrokers"), nil)
	if err != nil {
		beego.Error("error when creating Kafka SyncProducer", err)
		return err
	}
	defer func() {
		if errClose := kafkaProducer.Close(); errClose != nil {
			beego.Error("error when closing Kafka SyncProducer", errClose)
		}
	}()
	kafkaMessage := &sarama.ProducerMessage{
		Topic: beego.AppConfig.String("KafkaTopic"),
		Key:   nil,
		Value: sarama.StringEncoder(message.Title),
	}
	_, _, errSend := kafkaProducer.SendMessage(kafkaMessage)
	return errSend
}

func init() {
	DefaultMessageList = NewMessageManager()
}
