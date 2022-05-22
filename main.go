package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kafy11/gowsclient/client"
)

const WS_ADDRESS = "ws://localhost:8080?id=1"

func main() {
	test := flag.Bool("t", false, "Testar conexão")
	flag.Parse()

	if *test {
		client.New(WS_ADDRESS).Connect()
		fmt.Println("Sucesso ao conectar")
		os.Exit(0)
	}
}
