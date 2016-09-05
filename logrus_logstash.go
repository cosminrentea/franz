package main

import (
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"time"
	"github.com/smancke/guble/logformatter"
	log "github.com/Sirupsen/logrus"
)

type logrusLogstash struct {
	BeeLevel    int
	Level       string
	Env         string
	ServiceName string
}

// NewLogrusLogstash returns a LoggerInterface
func NewLogrusLogstash() logs.Logger {
	ll := &logrusLogstash{
		BeeLevel: logs.LevelError,
	}
	return ll
}

// Init with JSON config, like: {"Level": "debug", "Env": "int", "ServiceName": "franz"}
func (ll *logrusLogstash) Init(jsonconfig string) error {
	err := json.Unmarshal([]byte(jsonconfig), ll)
	if err != nil {
		return err
	}
	log.SetFormatter(&logformatter.LogstashFormatter{
		ServiceName: ll.ServiceName,
		Env: ll.Env,
	})
	switch ll.Level {
	case "debug":
		log.SetLevel(log.DebugLevel)
		ll.BeeLevel = logs.LevelDebug
	case "info":
		log.SetLevel(log.InfoLevel)
		ll.BeeLevel = logs.LevelInformational
	case "warn":
		log.SetLevel(log.WarnLevel)
		ll.BeeLevel = logs.LevelWarning
	case "error":
		fallthrough
	default:
		log.SetLevel(log.ErrorLevel)
		ll.BeeLevel = logs.LevelError
	}
	return nil
}

// WriteMsg will write the msg and level into ll.
func (ll *logrusLogstash) WriteMsg(when time.Time, msg string, level int) error {
	if level > ll.BeeLevel {
		return nil
	}
	switch level {
	case logs.LevelDebug:
		log.Debug(msg)
	case logs.LevelInformational:
		log.Info(msg)
	case logs.LevelWarning:
		log.Warn(msg)
	case logs.LevelError:
		log.Error(msg)
	}
	return nil
}

// Destroy is an empty method
func (ll *logrusLogstash) Destroy() {}

// Flush is an empty method
func (ll *logrusLogstash) Flush() {}
