package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sadmadrus/chessBox/cmd/game/gameserver"
)

var (
	address    = "localhost:8282"
	addressEnv = "LISTEN_ADDRESS"
)

func main() {
	http.HandleFunc("/", rootHandler)
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

// rootHandler отвечает за обработку запросов к сервису в целом.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		gameserver.GameHandler(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		fmt.Fprint(w, "The game server is online and working.")
	case http.MethodPost:
		gameserver.Creator(w, r)
	case http.MethodOptions:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		w.WriteHeader(http.StatusNoContent)
	default:
		w.Header().Set("Allow", "GET, POST, OPTIONS")
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
	}
}
