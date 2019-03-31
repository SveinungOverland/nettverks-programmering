package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PaesslerAG/gval"
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

	var result string

	request, err := readRequest(c)

	if err == nil {
		result = response(request)
	} else {
		result = "Closing connection"
	}

	c.Write([]byte(result))

	fmt.Println("Closing connection", c.RemoteAddr())
}

func readRequest(c net.Conn) (string, error) {
	reader := bufio.NewReader(c)

	for {
		ln, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Print("Error is", err)
			return "Error\n", nil
		}
		str := string(ln)
		fmt.Print(str)
		if strings.HasPrefix(str, "GET") {
			return handleGETRequest(str, *reader), nil
		} else if strings.HasPrefix(str, "END") {
			return "Ending connection", errors.New("")
		} else {
			value, err := gval.Evaluate(str, nil)
			if err != nil {
				c.Write([]byte(fmt.Sprint("Encountered error", err, "\n")))
			} else {
				c.Write([]byte(fmt.Sprint(value, "\n")))
			}
		}
	}
}

func handleGETRequest(query string, reader bufio.Reader) string {
	fmt.Println("Got here")
	var sb strings.Builder
	contentLength := 0
	sb.WriteString(query)
	for {
		ln, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Print("Error is", err)
			break
		}
		str := string(ln)
		fmt.Print(str)
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

func response(request string) string {
	var ulLines strings.Builder
	lines := strings.Split(request, "\n")
	for _, line := range lines {
		ulLines.WriteString(fmt.Sprint("<li>", line, "</li>"))
	}
	body := fmt.Sprintf(`<!DOCTYPE html><html><head><meta charset="UTF-8"><title>Hello world</title></head>
	<body>
		HELLO WORLD!
		<ul>
			%s
		</ul>
	</body></html>`, ulLines.String())
	var sb strings.Builder
	sb.WriteString("HTTP/1.1 200 OK\r\n")
	sb.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(body)))
	sb.WriteString("\r\n")
	sb.WriteString(body)

	return sb.String()
}
