package client

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"sync"

	"golang.org/x/net/websocket"
)

type WsClient struct {
	config *websocket.Config
	conn   *websocket.Conn
	m      sync.RWMutex
}

type WsClientParams struct {
	SSL      bool
	URL      string
	Headers  map[string]string
	User     string
	Password string
}

type Handler func(string)

func New(params *WsClientParams) (*WsClient, error) {
	var origin, endpoint string
	if params.SSL {
		origin = fmt.Sprintf("https://%s", params.URL)
		endpoint = fmt.Sprintf("wss://%s", params.URL)
	} else {
		origin = fmt.Sprintf("http://%s", params.URL)
		endpoint = fmt.Sprintf("ws://%s", params.URL)
	}

	if params.User != "" && params.Password != "" {
		auth := fmt.Sprintf("%s:%s", params.User, params.Password)
		auth_encoded := base64.StdEncoding.EncodeToString([]byte(auth))
		origin = fmt.Sprintf("%s?authorization=%s", origin, auth_encoded)
		endpoint = fmt.Sprintf("%s?authorization=%s", endpoint, auth_encoded)
	}

	config, err := websocket.NewConfig(endpoint, origin)
	if err != nil {
		return nil, err
	}

	if params.SSL {
		config.TlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	for key, value := range params.Headers {
		config.Header.Add(key, value)
	}

	return &WsClient{config: config}, nil
}

func (client *WsClient) ListenMessages(messageHandler Handler) error {
	messageChannel := make(chan string)
	errChannel := make(chan error)
	go readMessages(client, messageChannel, errChannel)

	for {
		select {
		case message := <-messageChannel:
			go messageHandler(message)

		case err := <-errChannel:
			return err
		}
	}
}

func (client *WsClient) Connect() error {
	//locka para bloquear leituras na variável enquanto estiver tentando conectar
	client.m.Lock()
	defer client.m.Unlock() //defer para desbloqueiar a variável no final da função

	ws, err := websocket.DialConfig(client.config)
	if err != nil {
		return err
	}

	client.conn = ws
	return nil
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
		return err
	}
	return nil
}
