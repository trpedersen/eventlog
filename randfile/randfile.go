package main

import (
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

	var blockNumber int
	for {
		//fmt.Print("\noffset length: ")
		fmt.Print("\nblock number: ")

		//offset, err = scanner.ReadLong()
		blockNumber, err = scanner.ReadInt()
		if err != nil {
			fmt.Println("Enter an int\n")
			continue
		}
		//length, err = scanner.ReadInt()
		//if err != nil {
		//	fmt.Println("Enter an int\n")
		//	continue
		//}
		bytes := make([]byte, BLOCK_LEN)
		var n int
		if n, err = file.ReadAt(bytes, int64(blockNumber)*BLOCK_LEN); (err != nil) && (n == 0) {
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
