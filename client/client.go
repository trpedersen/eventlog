package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	COUNT      = 100
)

//strings.NewReader(rand.RandStr(RANDSTRLEN, "alphanum"))
func run(port string, topic string, threads int, count int) {
	var wg sync.WaitGroup
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func(thread int) {
			defer wg.Done()
			for j := 1; j <= count; j++ {
				response, err := http.Post(fmt.Sprintf("http://localhost:%s/eventlogs/%s", port, topic), "text/plain", strings.NewReader(fmt.Sprintf("thd# %d evt# %d", thread, j)))
				if err != nil {
					log.Printf("%d POST error: %s\n", thread, err)
					return
				}
				response.Body.Close()
			}
		}(i)
	}
	//for i := 0; i < COUNT; i++ {
	//	wg.Add(1)
	//	go func() {
	//		defer wg.Done()
	//		for j := 0; j < COUNT; j++ {
	//			_, err := http.Get("http://localhost:8080/events/topic2")
	//			if err != nil {
	//				log.Printf("GET error: %s\n", err)
	//			}
	//		}
	//	}()
	//}
	wg.Wait()
}

func main() {
	port := flag.String("port", "8080", "port [8080]")
	topic := flag.String("topic", "topic", "topic [topic]")
	threads := flag.Int("threads", 1, "threads [1]")
	count := flag.Int("count", 1, "count [1]")
	flag.Parse()
	run(*port, *topic, *threads, *count)
}
