package client

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/net/websocket"
)

type MessageToSent struct {
	Action string `json:"action"`
	To     int    `json:"to"`
	Msg    string `json:"msg"`
}

type MessageReceived struct {
	Action string `json:"action"`
	From   int    `json:"from"`
	Msg    string `json:"msg"`
}

type MessageHandlerFunction func(*MessageReceived) *MessageToSent

func Start(address string, messageHandler MessageHandlerFunction) {
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
	ws, err := websocket.Dial(fmt.Sprintf("ws://%s?id=1", address), "", fmt.Sprintf("http://%s?id=1", address))

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
