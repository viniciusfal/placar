package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/viniciusfal/placar/internal/config"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao Carregar.env file")
	}

	_, err = config.Load()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(("tudo certo com a aplicação"))

}
