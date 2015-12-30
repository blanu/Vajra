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
		CheckError(acceptErr)

		stopDetecting := make(chan bool)
		go detectPorts(conn, stopDetecting)

    fmt.Println("Reading...")
		reader := bufio.NewReader(conn)
		var selectedPort uint16
		readErr := binary.Read(reader, binary.LittleEndian, &selectedPort)
		fmt.Println("Read.")
		stopDetecting <- true
		CheckError(readErr)

    if readErr == nil {
			fmt.Println("Selected port", selectedPort)
			capturePort(selectedPort)
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

func detectPorts(conn net.Conn, stopDetecting chan bool) {
	for {
		select {
			case <-stopDetecting:
				 return
			default:
				fmt.Println("detecting...")
				time.Sleep(1000 * time.Millisecond)
		}
	}
}

func capturePort(port uint16) {
	fmt.Println("Capturing port", port)
	for {
		fmt.Println("capturing...", port)
		time.Sleep(1000 * time.Millisecond)
  }
}
