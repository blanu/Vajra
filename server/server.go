package server

import (
	"fmt"
	"os"
	"time"

	//  "github.com/ugorji/go/codec"

	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/transport/tcp"
	//	"github.com/go-mangos/mangos/protocol/req"
	"github.com/go-mangos/mangos/protocol/rep"
)

type Responder func([]byte) []byte

type VajraServer struct {
	sock mangos.Socket
}

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func date() string {
	return time.Now().Format(time.ANSIC)
}

func Listen(url string) VajraServer {
	var sock mangos.Socket
	var err error

	if sock, err = rep.NewSocket(); err != nil {
		die("can't get new rep socket: %s", err)
	}

	sock.AddTransport(tcp.NewTransport())
	if err = sock.Listen(url); err != nil {
		die("can't listen on rep socket: %s", err.Error())
	}

	return VajraServer{
		sock: sock,
	}
}

func (self VajraServer) Accept(responder Responder) []byte {
	var err error
	var msg []byte
	var response []byte

	// Could also use sock.RecvMsg to get header
	msg, err = self.sock.Recv()
	//	fmt.Println("server received request: %s", string(msg))

	response = responder(msg)

	err = self.sock.Send(response)
	if err != nil {
		die("can't send reply: %s", err.Error())
	}

	return msg
}

func Ok(request []byte) []byte {
	return []byte("success")
}
