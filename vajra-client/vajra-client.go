package main

import (
	"os"
  "fmt"

	"github.com/blanu/Vajra/client"
	"github.com/blanu/Vajra/message"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func main() {
	var c client.VajraClient
	var msg []byte
	var resp []byte
	var err error
	var ports []uint16

	c = client.Connect("tcp://localhost:10001")

	msg, err = message.EncodeStart()
	CheckError(err)
	c.Request(msg)

	msg, err = message.EncodeGetPorts()
	CheckError(err)
	resp = c.Request(msg)

	ports, err = message.DecodeDetectPorts(resp)
	CheckError(err)
	fmt.Println(ports)

	msg, err = message.EncodeChoosePort(443)
	CheckError(err)
	c.Request(msg)

	msg, err = message.EncodeStop()
	CheckError(err)
	c.Request(msg)
}
