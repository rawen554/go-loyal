package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/rawen554/go-loyal/internal/app"
	"github.com/rawen554/go-loyal/internal/config"
	"github.com/rawen554/go-loyal/internal/store"
	"gorm.io/gorm/logger"
)

const (
	timeoutServerShutdown = time.Second * 5
	timeoutShutdown       = time.Second * 10
)

func main() {
	ctx, cancelCtx := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelCtx()

	config, err := config.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	storage, err := store.NewPostgresStore(ctx, config.DatabaseURI, logger.LogLevel(config.LogLevel))
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	defer func() {
		wg.Wait()
	}()

	wg.Add(1)
	go func() {
		defer log.Print("closed DB")
		defer wg.Done()
		<-ctx.Done()

		storage.Close()
	}()

	componentsErrs := make(chan error, 1)

	app := app.NewApp(config, storage)
	r := app.SetupRouter()

	srv := http.Server{
		Addr:    config.RunAddr,
		Handler: r,
	}

	go func(errs chan<- error) {
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			errs <- fmt.Errorf("run server has failed: %w", err)
		}
	}(componentsErrs)

	wg.Add(1)
	go func() {
		defer log.Print("server has been shutdown")
		defer wg.Done()
		<-ctx.Done()

		shutdownTimeoutCtx, cancelShutdownTimeoutCtx := context.WithTimeout(context.Background(), timeoutServerShutdown)
		defer cancelShutdownTimeoutCtx()
		if err := srv.Shutdown(shutdownTimeoutCtx); err != nil {
			log.Printf("an error occurred during server shutdown: %v", err)
		}
	}()

	select {
	case <-ctx.Done():
	case err := <-componentsErrs:
		log.Print(err)
		cancelCtx()
	}

	go func() {
		ctx, cancelCtx := context.WithTimeout(context.Background(), timeoutShutdown)
		defer cancelCtx()

		<-ctx.Done()
		log.Fatal("failed to gracefully shutdown the service")
	}()
}
