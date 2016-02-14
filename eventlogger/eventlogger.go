package eventlogger

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
	"errors"
	"github.com/trpedersen/eventlog/constants"
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

func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}

func (event Event) ToString() string {
	return fmt.Sprintf("[%s] [%s] [%s]", event.Time, event.Topic, string(event.Data[:clen(event.Data)]))
}

func (event Event) MarshalBinary() ([]byte, error) {
	buf := new( bytes.Buffer)
	var err error
	timeBytes, err := event.Time.MarshalBinary() // binary.Write(buf, binary.LittleEndian, event.Time.UnixNano()) // 8 bytes
	if err != nil {
		return nil, err
	}
	binary.Write(buf, binary.LittleEndian, timeBytes)
	var topicBytes = make([]byte, constants.EVENT_TOPIC_LEN, constants.EVENT_TOPIC_LEN)
	copy(topicBytes, []byte(event.Topic))
	err = binary.Write(buf, binary.LittleEndian, topicBytes)
	if err != nil {
		return nil, err
	}
	data := make([]byte, constants.EVENT_DATA_LEN, constants.EVENT_DATA_LEN)
	copy(data, event.Data)
	err = binary.Write(buf, binary.LittleEndian, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (event *Event) UnmarshalBinary(bytes []byte) error {
	if len(bytes) != constants.EVENT_BYTES_LEN {
		return errors.New("event.UnmarshalBinary: invalid length")
	}
	if err := event.Time.UnmarshalBinary(bytes[0:constants.EVENT_TIME_LEN]); err != nil {
		return err
	}
	event.Topic = string(bytes[constants.EVENT_TIME_LEN:(constants.EVENT_TIME_LEN+constants.EVENT_TOPIC_LEN)])
	event.Data = bytes[(constants.EVENT_TIME_LEN+constants.EVENT_TOPIC_LEN):]
	return nil
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
