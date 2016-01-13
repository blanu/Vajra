package main

import (
	"fmt"
	"net"
	"os"
	"bufio"
	"time"
	"encoding/binary"
	"github.com/deckarep/golang-set"
	"github.com/google/gopacket"
  "github.com/google/gopacket/pcap"
	"github.com/google/gopacket/layers"
)

type Connection struct {
	smaller layers.TCPPort
	larger layers.TCPPort
}

func NewConnection(first layers.TCPPort, second layers.TCPPort) Connection {
  if first <= second {
		return Connection{first, second}
	} else {
		return Connection{second, first}
	}
}

func (conn Connection) CheckPort(port layers.TCPPort) bool {
	return conn.smaller==port || conn.larger==port
}

//type Trace []gopacket.Packet

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

		captured := map[Connection][]layers.TCP {}

		reader := bufio.NewReader(conn)
		writer := bufio.NewWriter(conn)

		handle, pcapErr := pcap.OpenLive("en0", 1024, false, 30 * time.Second)
    if pcapErr != nil {
			handle.Close()
			continue
		}

		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		packetChannel := make(chan gopacket.Packet)
		go readPackets(packetSource, packetChannel)

		stopDetecting := make(chan bool)
		go detectPorts(writer, packetChannel, captured, stopDetecting)

    fmt.Println("Reading...")
		var selectedPort uint16
		readErr := binary.Read(reader, binary.LittleEndian, &selectedPort)
		fmt.Println("Read.")
		stopDetecting <- true

    if readErr == nil {
			fmt.Println("Selected port", selectedPort)
			discardUnusedPorts(layers.TCPPort(selectedPort), captured)
			capturePort(selectedPort, writer, packetChannel, captured)
		}

		handle.Close()
	}
}

/* A Simple function to verify error */
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func detectPorts(writer *bufio.Writer, packetChannel chan gopacket.Packet, captured map[Connection][]layers.TCP, stopDetecting chan bool) {
  ports := mapset.NewSet()

	for {
		select {
			case <-stopDetecting:
				 return
		  case packet := <-packetChannel:
				fmt.Println("detected.")
				fmt.Println(packet)
				fmt.Println(ports)

				// Let's see if the packet is TCP
		    tcpLayer := packet.Layer(layers.LayerTypeTCP)
		    if tcpLayer != nil {
		        fmt.Println("TCP layer detected.")
		        tcp, _ := tcpLayer.(*layers.TCP)

						if !ports.Contains(tcp.SrcPort) {
							binary.Write(writer, binary.LittleEndian, tcp.SrcPort)
							writer.Flush()
							ports.Add(tcp.SrcPort)
						}

						if !ports.Contains(tcp.DstPort) {
							binary.Write(writer, binary.LittleEndian, tcp.DstPort)
							writer.Flush()
							ports.Add(tcp.DstPort)
						}

						conn := NewConnection(tcp.SrcPort, tcp.DstPort)
						trace, ok := captured[conn]
						if !ok {
							trace = make([]layers.TCP, 1)
						}

						captured[conn]=append(trace, *tcp)
		    }
			default:
				fmt.Println("detecting...")
				time.Sleep(1000 * time.Millisecond)
		}
	}
}

func capturePort(port uint16, writer *bufio.Writer, packetChannel chan gopacket.Packet, captured map[Connection][]layers.TCP) {
	var count uint16 = uint16(len(captured))

	fmt.Println("Capturing port", port)

	for {
		fmt.Println("capturing...", port, count)
		select {
		  case packet := <-packetChannel:
				fmt.Println("detected.")
				fmt.Println(packet)

				// Let's see if the packet is TCP
		    tcpLayer := packet.Layer(layers.LayerTypeTCP)
		    if tcpLayer != nil {
		        fmt.Println("TCP layer detected.")
		        tcp, _ := tcpLayer.(*layers.TCP)

						conn := NewConnection(tcp.SrcPort, tcp.DstPort)
						if !conn.CheckPort(layers.TCPPort(port)) {
							continue
						}

						captured=recordPacket(tcp, captured)

						newCount := uint16(len(captured))
						if newCount > count {
							count = newCount
							writeErr := binary.Write(writer, binary.LittleEndian, count)
							if writeErr != nil {
								return
							}

							flushErr := writer.Flush()
							if flushErr != nil {
								return
							}
						}
		    }
			default:
				fmt.Println("capturing...")
				time.Sleep(1000 * time.Millisecond)
		}
	}
}

func readPackets(packetSource *gopacket.PacketSource, packetChannel chan gopacket.Packet) {
	fmt.Println("reading packets")
	for packet := range packetSource.Packets() {
		fmt.Println("readPacket")
		packetChannel <- packet
	}
	fmt.Println("done reading packets")
}

func discardUnusedPorts(port layers.TCPPort, captured map[Connection][]layers.TCP) {
	for conn := range captured {
		if !conn.CheckPort(port) {
			delete(captured, conn)
		}
	}
}

func recordPacket(tcp *layers.TCP, captured map[Connection][]layers.TCP) map[Connection][]layers.TCP {
	conn := NewConnection(tcp.SrcPort, tcp.DstPort)
	trace, ok := captured[conn]
	if !ok {
		trace = make([]layers.TCP, 1)
	}

	captured[conn]=append(trace, *tcp)

	return captured
}
