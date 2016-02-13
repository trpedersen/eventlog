package main

import (
	"flag"

	"os"

	"github.com/trpedersen/rand"
	"fmt"
)

const (
	COUNT      = 100
	RANDSTRLEN = 100
)

func run(pipe *os.File, topic string, threads int, count int) {
	n, err := pipe.Write([]byte(rand.RandStr(RANDSTRLEN, "alphanum")))
	fmt.Println(n, err)
	//var wg sync.WaitGroup
	//for i := 0; i < threads; i++ {
	//	wg.Add(1)
	//	go func(thread int) {
	//		defer wg.Done()
	//		for j := 0; j < count; j++ {
	//			_, err := pipe.WriteString(rand.RandStr(RANDSTRLEN, "alphanum")+ "\n")
	//			if err != nil {
	//				log.Printf("%d write error: %s\n", thread, err)
	//				return
	//			}
	//		}
	//	}(i)
	//}
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
	//wg.Wait()
}

func main() {
	///port := flag.String("port", "8080", "port [8080]")
//	pipe, err := os.OpenFile("//./pipe/eventlog", os.O_RDWR, os.ModeNamedPipe)
	pipe, err := os.OpenFile("eventlog", os.O_RDWR, os.ModeNamedPipe)
	defer pipe.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println(pipe)
	topic := flag.String("topic", "topic", "topic [topic]")
	threads := flag.Int("threads", 1, "threads [1]")
	count := flag.Int("count", 1, "count [1]")
	flag.Parse()
	run(pipe, *topic, *threads, *count)
}
