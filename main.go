package main

import (
	"github.com/kafy11/gowsclient/client"
)

func main() {
	client.Start("localhost:8080", messageHandler)
}

func messageHandler(received *client.MessageReceived) *client.MessageToSent {
	return &client.MessageToSent{
		Action: "teste",
		To:     received.From,
		Msg:    "Mensagem recebida",
	}
}
