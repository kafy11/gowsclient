package client

import (
	"fmt"
	"os"

	"golang.org/x/net/websocket"
)

type Handler func(string) string

func Start(address string, messageHandler Handler) {
	ws := connect(address)

	messageChannel := make(chan string)
	go readMessages(ws, messageChannel)

	for {
		select {
		case message := <-messageChannel:
			fmt.Println(`Message Received:`, message)

			err := websocket.JSON.Send(ws, messageHandler(message))
			if err != nil {
				fmt.Printf("Send failed: %s\n", err.Error())
				os.Exit(1)
			}
		}
	}
}

func Test(address string) {
	connect(address)
}

func connect(address string) *websocket.Conn {
	ws, err := websocket.Dial(address, "", fmt.Sprintf("http://%s", address))

	if err != nil {
		fmt.Printf("Falha ao conectar: %s\n", err.Error())
		os.Exit(1)
	}

	return ws
}

func readMessages(ws *websocket.Conn, incomingMessages chan string) {
	for {
		var message string
		err := websocket.Message.Receive(ws, &message)
		if err != nil {
			fmt.Printf("Error::: %s\n", err.Error())
			return
		}
		incomingMessages <- message
	}
}
