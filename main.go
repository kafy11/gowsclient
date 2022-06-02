package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kafy11/gowsclient/client"
)

const WS_ADDRESS = "ws://localhost:8080?id=1"

func main() {
	test := flag.Bool("t", false, "Testar conexão")
	flag.Parse()

	if *test {
		err := godotenv.Load(".env")

		if err != nil {
			log.Fatal("Erro no .env", "Error loading .env file")
		}

		//cria um client do websocket e connecta
		wsClient, err := client.New(&client.WsClientParams{
			SSL:      true,
			URL:      os.Getenv("WEBSOCKET_SERVER_URL"),
			User:     os.Getenv("WS_AUTH_USER"),
			Password: os.Getenv("WS_AUTH_PASS"),
		})

		if err != nil {
			log.Fatal("Falha ao criar o client ", err)
		}

		err = wsClient.Connect()

		if err != nil {
			log.Fatal("Falha ao conectar:", err)
		}

		fmt.Println("Sucesso ao conectar")
		os.Exit(0)
	}
}
