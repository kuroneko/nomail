// nomail.go
//
// A simple SMTP server that refuses to do anything
//

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"unicode"
)

const (
	SmtpNoServiceHere = 554
	SmtpBadSequenceOfCommands = 503
)

func sendError(conn net.Conn, code int, message string) {
	fmtmsg := fmt.Sprintf("%03d %s\n", code, message)
	conn.Write([]byte(fmtmsg))
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	tcpconn, ok := conn.(*net.TCPConn)
	if ok {
		tcpconn.SetLinger(-1)
		tcpconn.SetKeepAlive(true)
	}

	readbuf := bufio.NewReader(conn)
	sendError(conn, SmtpNoServiceHere, "This is not the SMTP Service you are looking for, move along.")
	for {
		cmd, err := readbuf.ReadString('\n')
		if err != nil {
			break
		}
		cmd = strings.TrimRightFunc(cmd, unicode.IsSpace)
		if cmd == "" {
			continue
		}
		if strings.ToLower(cmd) == "quit" {
			break
		}
		sendError(conn, SmtpBadSequenceOfCommands, "bad sequence of commands")
	}
}
	

func main() {
	listener, err := net.Listen("tcp", ":8025")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to bind to port 25: %s\n",
			err)
		os.Exit(1)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, 
				"error accepting connection: %s\n",
				err)
			break;
		}
		go handleConnection(conn)
	}

	os.Exit(0)
}

