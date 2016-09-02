package models

import (
	"testing"
)

func newMessageOrFatal(t *testing.T, title string) *Message {
	message, err := NewMessage("learn Go")
	if err != nil {
		t.Fatalf("new message: %v", err)
	}
	return message
}

func TestNewMessage(t *testing.T) {
	title := "learn Go"
	message := newMessageOrFatal(t, title)
	if message.Title != title {
		t.Errorf("expected title %q, got %q", title, message.Title)
	}
	if message.Done {
		t.Errorf("new message is done")
	}
}

func TestNewMessageEmptyTitle(t *testing.T) {
	_, err := NewMessage("")
	if err == nil {
		t.Errorf("expected 'empty title' error, got nil")
	}
}

func TestSaveMessageAndRetrieve(t *testing.T) {
	message := newMessageOrFatal(t, "learn Go")

	m := NewMessageManager()
	m.Save(message)

	all := m.All()
	if len(all) != 1 {
		t.Errorf("expected 1 message, got %v", len(all))
	}
	if *all[0] != *message {
		t.Errorf("expected %v, got %v", message, all[0])
	}
}

func TestSaveAndRetrieveTwoMessages(t *testing.T) {
	learnGo := newMessageOrFatal(t, "learn Go")
	learnTDD := newMessageOrFatal(t, "learn TDD")

	m := NewMessageManager()
	m.Save(learnGo)
	m.Save(learnTDD)

	all := m.All()
	if len(all) != 2 {
		t.Errorf("expected 2 messages, got %v", len(all))
	}
	if *all[0] != *learnGo && *all[1] != *learnGo {
		t.Errorf("missing message: %v", learnGo)
	}
	if *all[0] != *learnTDD && *all[1] != *learnTDD {
		t.Errorf("missing message: %v", learnTDD)
	}
}

func TestSaveModifyAndRetrieve(t *testing.T) {
	message := newMessageOrFatal(t, "learn Go")
	m := NewMessageManager()
	m.Save(message)

	message.Done = true
	if m.All()[0].Done {
		t.Errorf("saved message wasn't done")
	}
}

func TestSaveTwiceAndRetrieve(t *testing.T) {
	message := newMessageOrFatal(t, "learn Go")
	m := NewMessageManager()
	m.Save(message)
	m.Save(message)

	all := m.All()
	if len(all) != 1 {
		t.Errorf("expected 1 message, got %v", len(all))
	}
	if *all[0] != *message {
		t.Errorf("expected message %v, got %v", message, all[0])
	}
}

func TestSaveAndFind(t *testing.T) {
	message := newMessageOrFatal(t, "learn Go")
	m := NewMessageManager()
	m.Save(message)

	nt, ok := m.Find(message.ID)
	if !ok {
		t.Errorf("didn't find message")
	}
	if *message != *nt {
		t.Errorf("expected %v, got %v", message, nt)
	}
}
