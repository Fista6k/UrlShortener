package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Fista6k/Url-Shorterer.git/internal/adapter"
	controller "github.com/Fista6k/Url-Shorterer.git/internal/controller/http"
	"github.com/Fista6k/Url-Shorterer.git/internal/service"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found\n")
	}
}

func main() {

	port := flag.String("port", "8080", "where i run server")
	flag.Parse()

	storage, err := adapter.ConnToStorage()
	if err != nil {
		fmt.Println("something went wrong with storage init")
		fmt.Printf("err: %v", err)
		os.Exit(1)
	}

	service := service.NewShortererService(storage)
	r := controller.NewRouter(service)
	addr := ":" + *port
	if err = r.Router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
