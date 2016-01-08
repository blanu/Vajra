package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
	"time"
	"encoding/binary"
)

func main() {
	fmt.Println("Launching server...")
	ln, listenErr := net.Listen("tcp", ":10001")
	CheckError(listenErr)

	for {
		fmt.Println("Accepting connection...")
    conn, acceptErr := ln.Accept()
		fmt.Println("Accepted connection.")
		if acceptErr != nil {
			continue
		}

		reader := bufio.NewReader(conn)
		writer := bufio.NewWriter(conn)

		stopDetecting := make(chan bool)
		go detectPorts(writer, stopDetecting)

    fmt.Println("Reading...")
		var selectedPort uint16
		readErr := binary.Read(reader, binary.LittleEndian, &selectedPort)
		fmt.Println("Read.")
		stopDetecting <- true

    if readErr == nil {
			fmt.Println("Selected port", selectedPort)
			capturePort(selectedPort, writer)
		}
	}
}

/* A Simple function to verify error */
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func detectPorts(writer *bufio.Writer, stopDetecting chan bool) {
	var detectedPort uint16 = 80
	binary.Write(writer, binary.LittleEndian, detectedPort)

	detectedPort = 443
	binary.Write(writer, binary.LittleEndian, detectedPort)

	detectedPort = 1234
	binary.Write(writer, binary.LittleEndian, detectedPort)

	for {
		select {
			case <-stopDetecting:
				 return
			default:
				fmt.Println("detecting...")
				writer.Flush()
				time.Sleep(1000 * time.Millisecond)
		}
	}
}

func capturePort(port uint16, writer *bufio.Writer) {
	var count uint16 = 0

	fmt.Println("Capturing port", port)
	for {
		fmt.Println("capturing...", port, count)
		count = count + 1
		writeErr := binary.Write(writer, binary.LittleEndian, count)
		if writeErr != nil {
			return
		}

		flushErr := writer.Flush()
		if flushErr != nil {
			return
		}

		time.Sleep(1000 * time.Millisecond)
  }
}
