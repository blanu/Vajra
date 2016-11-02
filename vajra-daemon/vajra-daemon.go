package main

import (
	"os"
	"fmt"
	"time"
	"flag"

	"github.com/deckarep/golang-set"

	"github.com/google/gopacket"
  "github.com/google/gopacket/pcap"
	"github.com/google/gopacket/layers"

	"github.com/blanu/Vajra/server"
	"github.com/blanu/Vajra/message"

	"github.com/blanu/AdversaryLab-protocol/adversarylab"
)

type Connection struct {
	smaller layers.TCPPort
	larger layers.TCPPort
}

func NewConnection(packet *layers.TCP) Connection {
  if packet.SrcPort <= packet.DstPort {
		return Connection{packet.SrcPort, packet.DstPort}
	} else {
		return Connection{packet.DstPort, packet.SrcPort}
	}
}

func (conn Connection) CheckPort(port layers.TCPPort) bool {
	return conn.smaller==port || conn.larger==port
}

func main() {
	var lab adversarylab.Client
	var serve server.VajraServer
	var msg []byte
	var err error
	var captureName *string
	var datasetName *string

  captureName = flag.String("protocol", "", "Name of protocol being captured")
	datasetName = flag.String("dataset", "", "Name of dataset for captured data")
	flag.Parse()

	fmt.Println("Launching server...")

	lab = adversarylab.Connect("tcp://localhost:4567")
	serve = server.Listen("tcp://localhost:10001")

	fmt.Println("Accepting connection...")
	fmt.Println("Waiting for start.")
	msg = serve.Accept(server.Ok)
	err = message.DecodeStart(msg)
	CheckError(err)

	fmt.Println("Start.")

	captured := map[Connection]gopacket.Packet{}

	handle, pcapErr := pcap.OpenLive("en0", 1024, false, 30 * time.Second)
  if pcapErr != nil {
		handle.Close()
		os.Exit(1)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetChannel := make(chan gopacket.Packet)
	go readPackets(packetSource, packetChannel)

	stopDetecting := make(chan bool)
	ports := mapset.NewSet()
	go detectPorts(ports, packetChannel, captured, stopDetecting)

	serve.Accept(func (msg []byte) []byte {
		var resp []byte

		err = message.DecodeGetPorts(msg)
		CheckError(err)

    portObjs := ports.ToSlice()
		// The fact that Go requires this is evidence that Go types are not great.
		portNums := make([]uint16, len(portObjs))
		for i, arg := range portObjs {
			portNums[i] = arg.(uint16)
		}

		resp, err = message.EncodeDetectPorts(portNums)
		return resp
	})

  fmt.Println("Reading...")
	var selectedPort uint16
	msg = serve.Accept(server.Ok)
	selectedPort, err = message.DecodeChoosePort(msg)
	CheckError(err)

	fmt.Println("Read port.")
	stopDetecting <- true

	fmt.Println("Selected port", selectedPort)
	discardUnusedPorts(layers.TCPPort(selectedPort), captured)

	stopCapturing := make(chan bool)
	go capturePort(selectedPort, packetChannel, captured, stopCapturing)

	msg = serve.Accept(server.Ok)
	err = message.DecodeStop(msg)
	CheckError(err)
	stopDetecting <- true

	saveCaptured(lab, *captureName, *datasetName, captured)

	handle.Close()
}

/* A Simple function to verify error */
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func detectPorts(ports mapset.Set, packetChannel chan gopacket.Packet, captured map[Connection]gopacket.Packet, stopDetecting chan bool) {
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
							ports.Add(tcp.SrcPort)
						}

						if !ports.Contains(tcp.DstPort) {
							ports.Add(tcp.DstPort)
						}

            recordPacket(packet, captured)
		    }
			default:
				fmt.Println("detecting...")
				time.Sleep(1000 * time.Millisecond)
		}
	}
}

func capturePort(port uint16, packetChannel chan gopacket.Packet, captured map[Connection]gopacket.Packet, stopCapturing chan bool) {
	var count uint16 = uint16(len(captured))

	fmt.Println("Capturing port", port)

	for {
		fmt.Println("capturing...", port, count)
		select {
			case <-stopCapturing:
				 return
		  case packet := <-packetChannel:
				fmt.Println("detected.")
				fmt.Println(packet)

				// Let's see if the packet is TCP
		    tcpLayer := packet.Layer(layers.LayerTypeTCP)
		    if tcpLayer != nil {
		        fmt.Println("TCP layer detected.")
		        tcp, _ := tcpLayer.(*layers.TCP)

						conn := NewConnection(tcp)
						if !conn.CheckPort(layers.TCPPort(port)) {
							continue
						}

						recordPacket(packet, captured)

						newCount := uint16(len(captured))
						if newCount > count {
							count = newCount
							fmt.Println("%d packets", count)
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

func discardUnusedPorts(port layers.TCPPort, captured map[Connection]gopacket.Packet) {
	for conn := range captured {
		if !conn.CheckPort(port) {
			delete(captured, conn)
		}
	}
}

func recordPacket(packet gopacket.Packet, captured map[Connection]gopacket.Packet) {
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		fmt.Println("TCP layer detected.")
		tcp, _ := tcpLayer.(*layers.TCP)
		conn := NewConnection(tcp)
		_, ok := captured[conn]
		// Save first packet only for each new connection
		if !ok {
			captured[conn]=packet
		}
	}
}

func saveCaptured(lab adversarylab.Client, name string, dataset string, captured map[Connection]gopacket.Packet) {
	fmt.Println("Saving captured byte sequences...")

	for _, packet := range(captured) {
		if app := packet.ApplicationLayer(); app != nil {
			data := app.Payload()
			fmt.Println(data)
			lab.AddPacket(name, dataset, data)
	  }
	}
}
