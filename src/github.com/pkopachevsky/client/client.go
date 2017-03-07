package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

var (
	host = flag.String("h", "localhost", "Host")
	port = flag.Int("p", 0, "Port")
	execute = flag.String("e", "", "Execute command")
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 2 {
		fmt.Println("Hostname and port required")
		return
	}
	serverHost := flag.Arg(0)
	serverPort := flag.Arg(1)
	startClient(fmt.Sprintf("%s:%s", serverHost, serverPort))
}

func startClient(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Can't connect to server: %s\n", err)
		return
	}
	if len(*execute) > 0 {
		cmd := fmt.Sprintf("%s\n", *execute)
		conn.Write([]byte(cmd))
	}
	go io.Copy(connWriter{os.Stdout, os.Stderr}, conn)
	_, err = io.Copy(conn, os.Stdin)
	if err != nil {
		fmt.Printf("Connection error: %s\n", err)
	}
}

type connWriter struct {
	destOk io.Writer;
	destErr io.Writer;
}

func (cr connWriter) Write(p []byte) (n int, err error) {
	if(p[len(p)-1] == 0) {
		fmt.Println("Success!")
		return cr.destOk.Write(p[0:len(p)-1])
	} else {
		fmt.Println("Failure!")
		return cr.destErr.Write(p[0:len(p)-1])
	}
}
