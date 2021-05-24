package client

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kafy11/gosocket/message"
	"golang.org/x/net/websocket"
)

type MessageToSend message.Received
type MessageReceived message.ToSend

type Handler func(*MessageReceived) *MessageToSend

func Start(address string, messageHandler Handler) {
	ws := connect(address)

	messageChannel := make(chan *MessageReceived)
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
	ws, err := websocket.Dial(fmt.Sprintf("ws://%s", address), "", fmt.Sprintf("http://%s", address))

	if err != nil {
		fmt.Printf("Falha ao conectar: %s\n", err.Error())
		os.Exit(1)
	}

	return ws
}

func readMessages(ws *websocket.Conn, incomingMessages chan *MessageReceived) {
	for {
		var jsonInput string
		err := websocket.Message.Receive(ws, &jsonInput)
		if err != nil {
			fmt.Printf("Error::: %s\n", err.Error())
			return
		}
		message := new(MessageReceived)
		json.Unmarshal([]byte(jsonInput), &message)
		incomingMessages <- message
	}
}
