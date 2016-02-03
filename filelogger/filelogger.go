package filelogger

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/trpedersen/eventlog/eventlogger"
	"github.com/trpedersen/workqueue"
)

const (
	MAXWORKERS = 10
)

type fileEventLogger struct {
	topicFiles map[string]*os.File
	workqueue  workqueue.WorkQueue
}

func NewFileEventLogger() eventlogger.EventLogger {
	workqueue, err := workqueue.NewWorkQueue(MAXWORKERS)
	if err != nil {
		panic(err)
	}
	return &fileEventLogger{
		topicFiles: make(map[string]*os.File),
		workqueue:  workqueue,
	}
}

func (logger *fileEventLogger) Log(event eventlogger.Event) (err error) {
	job := &fileEventLogJob{
		logger: logger,
		event:  event,
	}
	return logger.workqueue.Enqueue(job)
}

func (logger *fileEventLogger) getTopicFile(topic string) (*os.File, error) {

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
	}
	return topicFile, err
}

type fileEventLogJob struct {
	logger *fileEventLogger
	event  eventlogger.Event
}

func (job *fileEventLogJob) Execute() error {
	topicFile, err := job.logger.getTopicFile(job.event.Topic)
	if err != nil {
		return err
	}
	if _, err = fmt.Fprintf(topicFile, "[%37s] %10s\n", job.event.Time, job.event.Data); err != nil {
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
