package main

import (
	"log"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/trpedersen/rand"
)

const (
	COUNT      = 100
	RANDSTRLEN = 100
)

func TestTopic(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < COUNT; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < COUNT; j++ {
				_, err := http.Post("http://localhost:8080/events/topic2", "text/plain", strings.NewReader(rand.RandStr(RANDSTRLEN, "alphanum")))
				if err != nil {
					log.Printf("POST error: %s\n", err)
				}
			}
		}()
	}
	for i := 0; i < COUNT; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < COUNT; j++ {
				_, err := http.Get("http://localhost:8080/events/topic2")
				if err != nil {
					log.Printf("GET error: %s\n", err)
				}
			}
		}()
	}
	wg.Wait()
}
