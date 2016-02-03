package main

import (
	"flag"
	"fmt"
	"github.com/trpedersen/io"
	"log"
	"os"
)

const (
	RECLEN = int64(51) // + 1L
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
	//var offset int64
	//var length int
	var n int
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
		bytes := make([]byte, RECLEN-1)
		n, err = file.ReadAt(bytes, int64(recordNumber)*RECLEN)
		fmt.Printf("%d bytes read, err: %s, bytes: %s", n, err, bytes)
		//fmt.Println(offset, length)
	}
}
