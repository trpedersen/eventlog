package main

import (
	"flag"
	"fmt"
	"github.com/trpedersen/io"
	"log"
	"os"
	"github.com/trpedersen/eventlog/eventlogger"
	"github.com/trpedersen/eventlog/constants"
)

const (
	RECLEN = int64(constants.EVENT_BYTES_LEN)
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

	var recordNumber int
	for {
		//fmt.Print("\noffset length: ")
		fmt.Print("\nrecord number: ")

		//offset, err = scanner.ReadLong()
		recordNumber, err = scanner.ReadInt()
		if err != nil {
			fmt.Println("Enter an int\n")
			continue
		}
		//length, err = scanner.ReadInt()
		//if err != nil {
		//	fmt.Println("Enter an int\n")
		//	continue
		//}
		bytes := make([]byte, RECLEN)
		if _ , err = file.ReadAt(bytes, int64(recordNumber)*RECLEN); err != nil {
			fmt.Println(err)
		} else {
			event := eventlogger.Event{}
			if err = event.UnmarshalBinary(bytes); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(event.ToString())
			}
		}
	}
}
