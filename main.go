package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kafy11/gosocket/log"
	"github.com/kafy11/gowsclient/client"
)

const WS_ADDRESS = "ws://localhost:8080?id=1"

func main() {
	test := flag.Bool("t", false, "Testar conexão")
	flag.Parse()

	if *test {
		client.Test(WS_ADDRESS)
		fmt.Println("Sucesso ao conectar")
		os.Exit(0)
	}
}

func messageHandler(received string) interface{} {
	log.Info("Mensagem recebida")
	return "Mensagem recebida"
}
