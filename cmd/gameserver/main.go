package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sadmadrus/chessBox/internal/gameserver"
)

var (
	address    = "localhost:8282"
	addressEnv = "LISTEN_ADDRESS"
)

func main() {
	stor := gameserver.NewMemoryStorage() // TODO: use a proper storage
	http.HandleFunc("/", gameserver.HandleRoot(stor))
	s := http.Server{
		Addr: address,
	}

	done := make(chan struct{}, 1)
	setupSignalHandling(done, &s)
	log.Printf("Running game service on %s.\n", address)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server exited with error: %v\n", err)
	}
	<-done
}

// init инициализирует конфигурацию из переменных окружения.
func init() {
	if value := os.Getenv(addressEnv); value != "" {
		address = value
	}
}

func setupSignalHandling(done chan<- struct{}, s *http.Server) {
	stopSig := make(chan os.Signal, 1)
	signal.Notify(stopSig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		<-stopSig
		log.Println("Shutting down...")
		if err := s.Shutdown(context.Background()); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
		close(done)
	}()
}
