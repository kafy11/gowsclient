package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kafy11/gosocket/log"
	"github.com/kafy11/gowsclient/client"
)

const WS_ADDRESS = "localhost:8080"

func main() {
	test := flag.Bool("t", false, "Testar conexão")
	flag.Parse()

	if *test {
		client.Test(WS_ADDRESS)
		fmt.Println("Sucesso ao conectar!")
		os.Exit(0)
	}

	client.Start(WS_ADDRESS, messageHandler)
}

func messageHandler(received *client.MessageReceived) *client.MessageToSent {
	log.Info("Mensagem recebida")
	return &client.MessageToSent{
		Action: "teste",
		To:     received.From,
		Msg:    "Mensagem recebida",
	}
}
