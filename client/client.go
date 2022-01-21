package client

import (
	"fmt"
	"os"

	"golang.org/x/net/websocket"
)

type WsClient struct {
	conn *websocket.Conn
}

type Handler func(string) interface{}

func Connect(address string) *WsClient {
	client := &WsClient{}
	client.conn = connect(address)

	return client
}

func (ws *WsClient) ListenMessages(messageHandler Handler) {
	messageChannel := make(chan string)
	go readMessages(ws.conn, messageChannel)

	for {
		select {
		case message := <-messageChannel:
			fmt.Println(`Message Received:`, message)

			msgToSend := messageHandler(message)
			if msgToSend == nil {
				continue
			}

			err := websocket.JSON.Send(ws.conn, msgToSend)
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

func (ws *WsClient) Send(message interface{}) error {
	err := websocket.JSON.Send(ws.conn, message)
	if err != nil {
		fmt.Printf("Send failed: %s\n", err.Error())
		return err
	}
	return nil
}
