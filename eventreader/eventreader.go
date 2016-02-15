package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/trpedersen/eventlog/eventlogger"
	"github.com/trpedersen/io"
	"log"
	"os"
)

const (
	BLOCK_LEN = int64(4096) //int64(constants.EVENT_BYTES_LEN)
)

func main() {
	scanner := io.NewScanner(os.Stdin)
	path := flag.String("path", "topic", "path+filename")
	flag.Parse()
	file, err := os.Open(*path) // todo: path for topic
	if err != nil {
		log.Printf("Error opening file, path: %s, err: %s\n", *path, err)
		return
	}

	var eventNumber int
	for {
		fmt.Print("\nevent number: ")

		eventNumber, err = scanner.ReadInt()
		if err != nil {
			fmt.Println("Enter an int\n")
			continue
		}

		var offset int64
		var recLen int32
		var event eventlogger.Event
		recLenBytes := make([]byte, 4)
		for i := 0; i <= eventNumber; i++ {
			if n, err := file.ReadAt(recLenBytes, offset); (err != nil) && (n == 0) {
				break // EOF
			}
			if err := binary.Read(bytes.NewReader(recLenBytes), binary.LittleEndian, &recLen); err != nil {
				panic(err)
			}
			if i == eventNumber {
				// we have our event
				eventBytes := make([]byte, recLen)
				if _, err := file.ReadAt(eventBytes, offset); err != nil {
					panic(err)
				}
				event.UnmarshalBinary(eventBytes)
				break
			} else {
				// seek to next event
				offset = offset + int64(recLen)
			}
		}
		fmt.Println(event.ToString())
	}
}
