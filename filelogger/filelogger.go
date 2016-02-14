package filelogger

import (
	"encoding/gob"
	"log"
	"os"
	"sync"

	"github.com/trpedersen/eventlog/eventlogger"
	"github.com/trpedersen/pubsub"
)

const (
	MAXWORKERS = 10
)

type fileLogger struct {
	files         map[string]*os.File
	encoders      map[string]*gob.Encoder
	hub           pubsub.Hub
	subscriptions map[string]pubsub.Subscription
	log           chan eventlogger.Event
	quit          chan struct{}
}

func NewFileLogger() eventlogger.EventLogger {
	hub, _ := pubsub.NewHub()
	logger := &fileLogger{
		files:         make(map[string]*os.File),
		encoders:      make(map[string]*gob.Encoder),
		hub:           hub,
		subscriptions: make(map[string]pubsub.Subscription),
		log:           make(chan eventlogger.Event),
	}
	go logger.run()
	return logger
}

func (logger *fileLogger) run() {
run:
	for {
		select {
		case <-logger.quit:
			break run
		case event := <-logger.log:
			logger.logEvent(event)
		}
	}
	close(logger.quit)
	return
}

func (logger *fileLogger) Log(event eventlogger.Event) (err error) {

	// todo: make this async
	logger.log <- event
	return nil
}

func (logger *fileLogger) Halt() {
	select {
	case <-logger.quit: // already closed
		return
	default:
		logger.quit <- struct{}{}
	}
	<-logger.quit
	return
}

func (logger *fileLogger) logEvent(event eventlogger.Event) (err error) {
	file, err := logger.getTopicFile(event.Topic)
	if err != nil {
		return err
	}
	//if _, err = fmt.Fprintf(file, "[%37s] %10s\n", event.Time, event.Data); err != nil {
	var bytes []byte
	bytes, err = event.MarshalBinary()
	if err != nil {
		log.Printf("Error marshalling event, err: %s", err)
		return err
	}
	if _, err = file.Write(bytes); err != nil {
		log.Printf("Error writing event, err: %s", err)
	}
	return err
}

func (logger *fileLogger) getTopicFile(topic string) (*os.File, error) {

	topicFile, exists := logger.files[topic]
	var err error
	if !exists {
		var mutex = &sync.Mutex{}
		mutex.Lock()
		defer mutex.Unlock()

		topicFile, err = os.OpenFile(topic, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666) // todo: path for topic
		if err != nil {
			log.Printf("Error opening file, topic: %s, err: %s\n", topic, err)
			return nil, err
		}
		logger.files[topic] = topicFile
		//logger.encoders[topic] = gob.NewDecoder(topicFile)
	}
	return topicFile, err
}

//func (this *EventLog) ReadEvents(topic string) ([]string, error) {
//
//	topicFile, err := os.Open(topic)
//	if err != nil {
//		log.Printf("Topic file error: %s\n", err)
//		return nil, err
//	}
//	s := bufio.NewScanner(topicFile)
//	events := make([]string, 0)
//	for s.Scan() {
//		events = append(events, s.Text())
//	}
//	return events, nil
//}

//func WriteEvent(topic string, event string) error {
//	events, exists := topics[topic]
//	if !exists {
//		events =  make([]string, 0)
//		topics[topic] = events
//	}
//	topics[topic] = append(events, event)
//	return nil
//}

//func ReadEvents(topic string) ([]string, error){
//	events, _ := topics[topic]
//	return events, nil
//}
