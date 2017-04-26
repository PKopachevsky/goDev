package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var (
	execute = flag.String("e", "", "Execute command")
	upload  = flag.String("u", "", "Upload file")
	host = flag.String("h", "localhost", "Host")
	port = flag.Int("p", 8888, "Port")
)

func main() {
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *host, *port)
	err := startClient(addr, *execute, *upload)
	if err != nil {
		fmt.Printf("Error: %s.\n", err.Error())
	}
}

func startClient(addr string, execute string, upload string) error{
	fmt.Printf("Conneciting to %s\n", addr)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("Can't connect to server: %s\n", err)
		return err
	}
	fmt.Printf("Connected.\n")

	if len(execute) > 0 && len(upload) > 0 {
		return fmt.Errorf("Can't execute command and upload file simultaneously")
	}
	switch {
		case len(execute) > 0:
			err = commandClient(execute, conn)
		case len(upload) > 0:
			err = fileClient(upload, conn)
		default:
			err = defaultClient(conn)
		}

	return err
}

func fileClient(filename string, conn io.WriteCloser) error {
	fmt.Println("Starting file server")

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	return uploadFile(stat.Name(), conn, file)
}

func uploadFile(name string, conn io.WriteCloser, file io.ReadCloser) error {
	fmt.Printf("Uploading file %s", name)
	_, err := io.WriteString(conn, fmt.Sprintf("%s_%d\n", name, time.Now()))
	if err != nil {
		return err
	}
	_, err = io.Copy(conn, file)
	if err != nil {
		return err
	}
	err = conn.Close()
	if err != nil {
		return err
	}
	return file.Close();
}

func commandClient(command string, conn io.ReadWriteCloser) error {
	_, err := io.WriteString(conn, fmt.Sprintf("%s\n", command))
	if err != nil {
		return err
	}
	return defaultClient(conn)
}

func defaultClient(conn io.ReadWriteCloser) error {
	go io.Copy(connWriter{os.Stdout, os.Stderr}, conn)
	_, err := io.Copy(conn, os.Stdin)
	if err != nil {
		return err
	}
	return conn.Close()
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
