package main

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/trpedersen/eventlog/eventlogger"
	"github.com/trpedersen/eventlog/filelogger"

	"bufio"
	"fmt"
	"os"
)

var eventLog = filelogger.NewFileEventLogger()

//func eventsHandlerGET( w http.ResponseWriter, r *http.Request){
//	urlParams := mux.Vars(r)
//	topic := urlParams["topic"]
//	events, _ := eventLog.ReadEvents(topic)
//	fmt.Fprint(w, events)
//}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	topic := urlParams["topic"]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic("ioutil.ReadAll")
	}
	eventLog.Log(eventlogger.NewEvent(time.Now(), topic, body))
	w.WriteHeader(http.StatusOK)
}

func Run() {

	pipe, err := os.OpenFile("eventlog2", os.O_RDWR|os.O_CREATE, os.ModeNamedPipe)
	println(pipe, err)
	//	pipe, err := os.OpenFile("\\\\.\\pipes\\eventlog", os.O_RDWR|os.O_CREATE, os.ModeNamedPipe)
	if err != nil {
		panic(err)
	}
	defer pipe.Close()

	reader := bufio.NewReader(pipe)
	data := make([]byte, 10)
	fmt.Println(data)
	n, err := reader.Read(data)
	fmt.Println(n, data)

}

func main() {
	//port := flag.String("port", "8080", "port [8080]")
	//flag.Parse()
	Run()
}
