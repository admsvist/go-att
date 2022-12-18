package app

import (
	"context"
	"fmt"
	"github.com/admsvist/go-att/api/handler"
	"github.com/admsvist/go-att/api/router"
	"github.com/admsvist/go-att/internal/app/repository"
	"github.com/admsvist/go-att/internal/app/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	s      *storage.Storage
	r      *repository.Repository
	h      *handler.Handler
	router *router.Router
}

func New() *App {
	a := &App{}
	a.s = storage.New()
	a.r = repository.New(a.s)
	a.h = handler.New(a.r)
	a.router = router.New(a.h)

	return a
}

func (a App) Run() error {
	fmt.Println("server running")

	// Load storage
	a.s.Read()

	// The HTTP Server
	server := &http.Server{Addr: ":8080", Handler: a.router.GetRouter()}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Save storage
		a.s.Write()

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()

	return nil
}
