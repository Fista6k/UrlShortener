package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Fista6k/Url-Shorterer.git/internal/adapter"
	controller "github.com/Fista6k/Url-Shorterer.git/internal/controller/http"
	"github.com/Fista6k/Url-Shorterer.git/internal/service"
)

func main() {
	port := flag.String("port", "8080", "where i run server")

	storage, err := adapter.ConnToStorage()
	if err != nil {
		fmt.Println("something went wrong with storage init")
		os.Exit(1)
	}

	service := service.NewShortererService(storage)
	r := controller.NewRouter(service)
	addr := ":" + *port
	if err = r.Router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
