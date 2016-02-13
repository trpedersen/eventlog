package eventlogger

import (
	"fmt"
	"time"
	"bytes"
	"encoding/gob"
)

//type Event interface {
//	ToString() string
//	GetTime() time.Time
//	GetTopic() string
//	GetData() []byte
//}

type Event struct {
	Time  time.Time
	Topic string
	Data  []byte
}

func NewEvent(time time.Time, topic string, data []byte) Event {
	return Event{
		Time:  time,
		Topic: topic,
		Data:  data,
	}
}

func (event Event) ToString() string {
	return fmt.Sprintf("[%s] %s", event.Time, event.Data)
}

func (event Event) Bytes() []byte {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	encoder.Encode(event)
	return buf.Bytes()
}

//func (event *basicEvent) GetTime() time.Time {
//	return event.time
//}

//
//func (event *basicEvent) GetTopic() string {
//	return event.topic
//}
//
//func (event *basicEvent) GetData() []byte {
//	return event.data
//}

type EventLogger interface {
	Log(event Event) error
	Halt()
}
