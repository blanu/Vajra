package main

import (
	"fmt"
	"net"
	"bufio"
	"testing"
	"encoding/binary"
)

func TestSelectedPort(t *testing.T) {
	fmt.Println("Connecting...")
	conn, _ := net.Dial("tcp", "127.0.0.1:10001")

	writer := bufio.NewWriter(conn)

	var selectedPort uint16 = 80
	binary.Write(writer, binary.LittleEndian, selectedPort)
	writer.Flush()
}
