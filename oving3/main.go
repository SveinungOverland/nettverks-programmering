package main

import (
	"fmt"
	"net"
	"time"

	"github.com/PaesslerAG/gval"
)

const maxBufferSize = 1024

func main() {
	address := ":8080"

	server(address)
}

func server(address string) (err error) {
	pc, err := net.ListenPacket("udp", address)
	if err != nil {
		return
	}
	defer pc.Close()

	fmt.Println("Listening for UDP on port: ", address)



	doneChan := make(chan error, 1)
	buffer := make([]byte, maxBufferSize)

	go func() {
		for {
			n, addr, err := pc.ReadFrom(buffer)
			if err != nil {
				doneChan <- err
				return
			}

			fmt.Printf("packet-received: bytes=%d from=%s\n", n, addr.String())

			deadline := time.Now().Add(15 * time.Second)
			err = pc.SetWriteDeadline(deadline)
			if err != nil {
				doneChan <- err
				return
			}

			str := string(buffer[:n])

			value, err := gval.Evaluate(str, nil)
			if err != nil {
				n, err = pc.WriteTo([]byte("Could not evaluate\n"), addr)
				if err != nil {
					doneChan <- err
					return
				}
			}
			n, err = pc.WriteTo([]byte(fmt.Sprint(value, "\n")), addr)
			if err != nil {
				doneChan <- err
				return
			}

			fmt.Printf("packet-written: bytes=%d to=%s\n", n, addr.String())
		}
	}()

	select {
	case err = <-doneChan:
	}

	return
}
