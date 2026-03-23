package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
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

	server := &http.Server{
		Handler: r.Router,
		Addr:    addr,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()

	stop()
	log.Println("shutting down gracefully, press Ctrl+C to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
