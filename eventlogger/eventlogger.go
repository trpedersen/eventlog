package eventlogger

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
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
	buf := new(bytes.Buffer)
	var err error
	timeBytes, err := event.Time.MarshalBinary()
	if err != nil {
		return nil, err
	}
	timeBytesLen := int32(len(timeBytes))

	topicBytes := []byte(event.Topic)
	topicBytesLen := int32(len(topicBytes))

	dataBytesLen := int32(len(event.Data))

	recLen := int32((4 * 4) + timeBytesLen + topicBytesLen + dataBytesLen)

	// record length
	if err = binary.Write(buf, binary.LittleEndian, recLen); err != nil {
		return nil, err
	}

	// time offset
	if err = binary.Write(buf, binary.LittleEndian, int32(16)); err != nil {
		return nil, err
	}

	// topic offset
	if err = binary.Write(buf, binary.LittleEndian, int32(16+timeBytesLen)); err != nil {
		return nil, err
	}

	// data offset
	if err = binary.Write(buf, binary.LittleEndian, int32(16+timeBytesLen+topicBytesLen)); err != nil {
		return nil, err
	}

	if err = binary.Write(buf, binary.LittleEndian, timeBytes); err != nil {
		return nil, err
	}
	if err = binary.Write(buf, binary.LittleEndian, topicBytes); err != nil {
		return nil, err
	}
	if err = binary.Write(buf, binary.LittleEndian, event.Data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (event *Event) UnmarshalBinary(data []byte) error {

	var recLen int32
	var err error
	if err = binary.Read(bytes.NewReader(data[0:4]), binary.LittleEndian, &recLen); err != nil {
		return err
	}
	var timeOffset int32
	if err = binary.Read(bytes.NewReader(data[4:8]), binary.LittleEndian, &timeOffset); err != nil {
		return err
	}
	var topicOffset int32
	if err = binary.Read(bytes.NewReader(data[8:12]), binary.LittleEndian, &topicOffset); err != nil {
		return err
	}
	var dataOffset int32
	if err = binary.Read(bytes.NewReader(data[12:16]), binary.LittleEndian, &dataOffset); err != nil {
		return err
	}
	//fmt.Printf("recLen: %d, timeOffset: %d, topicOffset: %d, dataOffset: %d, data: %s\n", recLen, timeOffset, topicOffset, dataOffset, data)

	if err = event.Time.UnmarshalBinary(data[timeOffset:topicOffset]); err != nil {
		return err
	}
	//fmt.Printf("event.Time: %s\n", event.Time)
	event.Topic = string(data[topicOffset:dataOffset])
	//fmt.Printf("event.Topic: %s\n", event.Topic)
	event.Data = data[dataOffset:recLen]
	//fmt.Printf("event.Data: %s\n", event.Data)
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
