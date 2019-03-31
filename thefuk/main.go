package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	PORT := ":" + arguments[1]
	listener, err := net.Listen("tcp4", PORT)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer listener.Close()
	rand.Seed(time.Now().Unix())

	fmt.Println("Listening on port", PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	defer c.Close()

	request := readRequest(c)

	fmt.Println(request)

	result := response()

	c.Write([]byte(result))
	fmt.Println("Closing connection", c.RemoteAddr())
}

func readRequest(c net.Conn) string {
	reader := bufio.NewReader(c)
	var sb strings.Builder
	contentLength := 0
	for {
		ln, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println(err)
		}
		str := string(ln)
		if str == "\r\n" {
			break
		} else if strings.HasPrefix(str, "content-length") {
			i, err := strconv.ParseInt(strings.TrimSpace(strings.Split(str, ":")[1]), 10, 32)
			if err != nil {
				panic("Unusable content-length")
			}
			contentLength = int(i)
		}
		sb.WriteString(str)
	}

	if contentLength != 0 {
		fmt.Println("Reading content of length", contentLength)
		buffer := make([]byte, contentLength)

		reader.Read(buffer)

		sb.WriteString(string(buffer))
	}

	return sb.String()
}

func response() string {
	body := `<!DOCTYPE html><html><head><meta charset="UTF-8"><title>Hello world</title></head><body>HELLO WORLD!</body></html>`
	var sb strings.Builder
	sb.WriteString("HTTP/1.1 200 OK\r\n")
	sb.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(body)))
	sb.WriteString("\r\n")
	sb.WriteString(body)

	return sb.String()
}
