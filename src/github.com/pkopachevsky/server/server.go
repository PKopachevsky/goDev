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
	bufReader := bufio.NewReader(conn)
	line, err := bufReader.ReadString('\n')
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

	reader := connReader{conn}
	okWriter := connWriter{conn}
	errWriter := errorWriter{conn}

	go io.Copy(stdin, reader)
	go io.Copy(okWriter, stdout)
	go io.Copy(errWriter, stderr)

	return cmd.Run()
}

type connReader struct {
	Conn net.Conn;
}

func (c connReader) Read(p []byte) (n int, err error) {
	n, err = c.Conn.Read(p);
	fmt.Printf("Input:\t%s", string(p));
	return n, err;
}

type connWriter struct {
	Conn net.Conn;
}

func (c connWriter) Write(p []byte) (n int, err error) {
	fmt.Printf("Output:\n%s\n", string(p));
	okResult := append(p, 0)
	return c.Conn.Write(okResult);
}

type errorWriter struct {
	Conn net.Conn;
}

func (c errorWriter) Write(p []byte) (n int, err error) {
	fmt.Printf("Error:%s\n", string(p));
	errorResult := append(p, 1)
	return c.Conn.Write(errorResult);
}
