package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/trpedersen/eventlog/eventlogger"
	"github.com/trpedersen/eventlog/filelogger"
)

var eventLog = filelogger.NewFileEventLogger()

//func eventsHandlerGET( w http.ResponseWriter, r *http.Request){
//	urlParams := mux.Vars(r)
//	topic := urlParams["topic"]
//	events, _ := eventLog.ReadEvents(topic)
//	fmt.Fprint(w, events)
//}

func eventsHandlerPOST(w http.ResponseWriter, r *http.Request) {
	urlParams := mux.Vars(r)
	topic := urlParams["topic"]
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic("ioutil.ReadAll")
	}
	eventLog.Log(eventlogger.NewEvent(time.Now(), topic, body))
	w.WriteHeader(http.StatusOK)
}

func Run(port *string) {

	gorillaRoute := mux.NewRouter()
	gorillaRoute.HandleFunc("/eventlogs/{topic}", eventsHandlerPOST).Methods("POST")
	//gorillaRoute.HandleFunc("/eventlogs/{topic}", eventsHandlerGET).Methods("GET")
	http.Handle("/", gorillaRoute)

	http.ListenAndServe(":"+*port, nil)
}

func main() {
	port := flag.String("port", "8080", "port [8080]")
	flag.Parse()
	Run(port)
}
