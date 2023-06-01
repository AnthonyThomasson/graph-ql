package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type server struct {
	router chi.Router
}

func NewServer(db *gorm.DB) *server {
	return &server{
		router: getRouter(db),
	}
}

func (s *server) Router() chi.Router {
	return s.router
}

func (s *server) Start(port string) {
	server := &http.Server{
		Addr:    port,
		Handler: s.router,
	}

	go func() {
		fmt.Println("Listening on port ", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("Server gracefully stopped")
}
