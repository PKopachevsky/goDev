package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"bufio"
	"os/exec"
	"strings"
)

var (
	host = flag.String("h", "localhost", "Host")
	port = flag.Int("p", 0, "Port")
)

func main() {
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		panic(err)
	}

	log.Printf("Listening for connections on %s", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection from client: %s", err)
		} else {
			go processClient(conn)
		}
	}
}

func processClient(conn net.Conn) {
	err := launchCommand(conn)
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}
	_, err2 := io.Copy(os.Stdin, conn)
	if err2 != nil {
		log.Println(err)
	}
	conn.Close()
}
func launchCommand(conn net.Conn) error {
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Printf("Command: %s\n", line)
	cmd := exec.Command(strings.TrimSpace(line))
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	go io.Copy(stdin, conn)
	go io.Copy(conn, stdout)
	go io.Copy(conn, stderr)

	return cmd.Run()
}
