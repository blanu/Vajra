package client

import (
	"os"
	"fmt"
	"time"

	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/transport/tcp"
	"github.com/go-mangos/mangos/protocol/req"
	//	"github.com/go-mangos/mangos/protocol/rep"
)

type VajraClient struct {
  sock mangos.Socket
}

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func date() string {
	return time.Now().Format(time.ANSIC)
}

func Connect(url string) VajraClient {
	var sock mangos.Socket
	var err error

	if sock, err = req.NewSocket(); err != nil {
		die("can't get new req socket: %s", err.Error())
	}

	sock.AddTransport(tcp.NewTransport())

	if err = sock.Dial(url); err != nil {
		die("can't dial on req socket: %s", err.Error())
	}

	return VajraClient {
    sock: sock,
  }
}

func (self VajraClient) Request(data []byte) []byte {
	var err error
	var msg []byte

	if err = self.sock.Send(data); err != nil {
		die("can't send message on push socket: %s", err.Error())
	}

	if msg, err = self.sock.Recv(); err != nil {
		die("can't receive date: %s", err.Error())
	}

	return msg
}
