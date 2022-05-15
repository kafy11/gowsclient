package client

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/net/websocket"
)

type WsClient struct {
	conn *websocket.Conn
}

type Handler func(string)

func Connect(address string) *WsClient {
	client := &WsClient{}
	client.conn = connect(address)

	return client
}

func (ws *WsClient) ListenMessages(messageHandler Handler) error {
	messageChannel := make(chan string)
	errChannel := make(chan error)
	go readMessages(ws.conn, messageChannel, errChannel)

	for {
		select {
		case message := <-messageChannel:
			fmt.Println(`Message Received:`, message)
			go messageHandler(message)

		case err := <-errChannel:
			return err
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

func readMessages(ws *websocket.Conn, incomingMessages chan string, errChannel chan error) {
	for {
		var message string
		err := websocket.Message.Receive(ws, &message)
		if err != nil {
			errChannel <- err
			return
		}
		incomingMessages <- message
	}
}

func (ws *WsClient) Send(message interface{}) error {
	if ws.conn == nil {
		return errors.New("failed to send message because websocket is not connected")
	}

	err := websocket.JSON.Send(ws.conn, message)
	if err != nil {
		fmt.Printf("Send failed: %s\n", err.Error())
		return err
	}
	return nil
}
