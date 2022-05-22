package client

import (
	"fmt"
	"os"
	"sync"

	"golang.org/x/net/websocket"
)

type WsClient struct {
	address string
	conn    *websocket.Conn
	m       sync.RWMutex
}

type Handler func(string)

func New(address string) *WsClient {
	return &WsClient{address: address}
}

func (client *WsClient) ListenMessages(messageHandler Handler) error {
	messageChannel := make(chan string)
	errChannel := make(chan error)
	go readMessages(client, messageChannel, errChannel)

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

func (client *WsClient) Connect() {
	//locka para bloquear leituras na variável enquanto estiver tentando conectar
	client.m.Lock()
	defer client.m.Unlock() //defer para desbloqueiar a variável no final da função

	ws, err := websocket.Dial(client.address, "", fmt.Sprintf("http://%s", client.address))

	if err != nil {
		fmt.Printf("Falha ao conectar: %s\n", err.Error())
		os.Exit(1)
	}

	client.conn = ws
}

func readMessages(client *WsClient, incomingMessages chan string, errChannel chan error) {
	for {
		//cria um lock de leitura
		client.m.RLock()

		var message string
		err := websocket.Message.Receive(client.conn, &message)
		if err != nil {
			errChannel <- err

			//Desbloqueia a o lock de leitura
			client.m.RUnlock()
			return
		}
		incomingMessages <- message

		//Desbloqueia a o lock de leitura
		client.m.RUnlock()
	}
}

func (client *WsClient) Send(message interface{}) error {
	//cria um lock de leitura
	client.m.RLock()
	defer client.m.RUnlock() //defer para desbloquear o lock no final da função

	err := websocket.JSON.Send(client.conn, message)
	if err != nil {
		fmt.Printf("Send failed: %s\n", err.Error())
		return err
	}
	return nil
}
