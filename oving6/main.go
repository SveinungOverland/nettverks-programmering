package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	PORT := ":" + arguments[1]

	listener, _ := net.Listen("tcp4", PORT)
	defer listener.Close()

	fmt.Println("Listening on port", PORT)

	for {
		conn, _ := listener.Accept()
		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	//defer c.Close()

	request := make(map[string]string)
	readRequest(c, request)

	fmt.Println("Request")
	pprintMap(request)

	isWsRequest := checkIfWsRequest(request)
	if !isWsRequest {
		panic("Not a websocket request, go use another server... pleb")
	}

	encodedKey := encodeKey(request["Sec-WebSocket-Key"])

	response := createHandshakeResponse(encodedKey)

	fmt.Println("Response")
	fmt.Println(response)

	c.Write([]byte(response))
}

func readRequest(c net.Conn, r map[string]string) error {
	reader := bufio.NewReader(c)

	ln, _ := reader.ReadBytes('\n')
	str := string(ln)

	if !strings.HasPrefix(str, "GET") {
		return errors.New("Something went wrong")
	}
	fmt.Println("FIRST LINE", str)
	r["Method"] = "GET"
	r["Request-URI"] = strings.TrimSuffix(str, "\r\n")

	for {
		ln, _ := reader.ReadBytes('\n')
		str := string(ln)
		if str == "\r\n" {
			break
		}

		parts := strings.SplitN(str, ":", 2)
		if len(parts) != 2 {
			return errors.New("Map parsing failed")
		}

		r[parts[0]] = parts[1][1 : len(parts[1])-2]

		// strings.TrimSuffix(parts[1], "\r\n") // This is the good way of removing newlines at the end of a line
		// = parts[1][:len(parts[1])-2]  is the cool kids way
	}

	return nil
}

func pprintMap(m map[string]string) {
	b, _ := json.MarshalIndent(m, "", "  ")
	fmt.Print(string(b))
}

func encodeKey(key string) string {
	h := sha1.New()
	io.WriteString(h, key+"258EAFA5-E914-47DA-95CA-C5AB0DC85B11")
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func checkIfWsRequest(r map[string]string) bool {
	mustHave := map[string]string{
		"Upgrade":               "websocket",
		"Connection":            "Upgrade",
		"Sec-WebSocket-Version": "13",
	}

	for k, v := range mustHave {
		if r[k] != v {
			return false
		}
	}

	return true
}

func createHandshakeResponse(encodedKey string) string {
	resp := map[string]string{
		"Upgrade":              "websocket",
		"Connection":           "Upgrade",
		"Sec-WebSocket-Accept": encodedKey,
	}
	return "HTTP/1.1 101 Switching Protocols\r\n" + mapToString(resp) + "\r\n"
}

func mapToString(m map[string]string) string {
	buffer := new(bytes.Buffer)
	for k, v := range m {
		fmt.Fprintf(buffer, "%s: %s\r\n", k, v)
	}
	return buffer.String()
}
