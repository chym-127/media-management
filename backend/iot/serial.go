package iot

import (
	"chym/stream/backend/api"
	"chym/stream/backend/protocols"
	"encoding/json"
	"fmt"
	"strings"

	"go.bug.st/serial"
)

func InitSerial() {
	// ports, err := serial.GetPortsList()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// if len(ports) == 0 {
	// 	fmt.Println("No serial ports found!")
	// }
	// for _, port := range ports {

	// 	fmt.Println("Found port: %v\n", port)
	// }

	mode := &serial.Mode{
		BaudRate: 9600,
		DataBits: 8,
		StopBits: 1,
	}
	port, err := serial.Open("COM5", mode)
	if err != nil {
		fmt.Println(err)
		return
	}
	go func() {
		for {
			buff := make([]byte, 200)
			n, err := port.Read(buff)
			if err != nil {
				fmt.Println(err)
				break
			}
			if n == 0 {
				fmt.Println("\nEOF")
			}
			key := string(buff[:n])
			key = strings.ReplaceAll(strings.ReplaceAll(key, "\n", ""), "\r", "")
			if len(key) > 2 {
				msg := protocols.EventMsg{
					Name: "IRREMOTE",
					Data: key,
				}
				jsonByte, _ := json.Marshal(msg)
				api.MsgStream.Message <- string(jsonByte[:])
			}
		}
	}()
}
