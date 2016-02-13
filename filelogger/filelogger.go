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
	files map[string]*os.File
	encoders map[string]*gob.Encoder
	hub  pubsub.Hub
	subscriptions map[string]pubsub.Subscription
	log chan eventlogger.Event
	quit chan struct{}
}

func NewFileLogger() eventlogger.EventLogger {
	logger := &fileLogger{
		files: make(map[string]*os.File),
		encoders: make(map[string]*gob.Encoder),
		hub:  pubsub.NewHub(),
		subscriptions: make(map[string]pubsub.Subscription),
		log: make(chan eventlogger.Event),
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
			logger
		}
	}
	close(logger.quit)
	return
}

func (logger *fileLogger) Log(event eventlogger.Event) (err error) {

	logger.log <- event
	return nil
	//msg := pubsub.NewMsg(event.Topic, event.Bytes())
	//logger.hub.Publish()
}

func (logger *fileLogger) Halt(){
	select {
	case <-logger.quit: // already closed
		return nil
	default:
		logger.quit <- struct{}{}
	}
	<-logger.quit
	return nil
}

func (logger *fileLogger) getTopicFile(topic string) (*os.File, error) {

	topicFile, exists := logger.topicFiles[topic]
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
		logger.topicFiles[topic] = topicFile
		logger.encoders[topic] = gob.NewDecoder(topicFile)
	}
	return topicFile, err
}

type fileEventLogJob struct {
	logger *fileLogger
	event  eventlogger.Event
}

func (job *fileEventLogJob) Execute() error {
	topicFile, err := job.logger.getTopicFile(job.event.Topic)
	if err != nil {
		return err
	}
//	if _, err = fmt.Fprintf(topicFile, "[%37s] %10s\n", job.event.Time, job.event.Data); err != nil {
	if err = job.logger.encoders[.Write(topicFile, binary.LittleEndian, job.event); err != nil {
		log.Printf("Error writing event, err: %s", err)
	}
	return err
}

func (job *fileEventLogJob) Report(err error) {
	// ignore it
	return
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
