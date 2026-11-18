package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/arinamklvch/go-courses/docs"
	"github.com/arinamklvch/go-courses/internal/adapter"
	"github.com/arinamklvch/go-courses/internal/config"
	grpchandler "github.com/arinamklvch/go-courses/internal/handler/grpc"
	httphandler "github.com/arinamklvch/go-courses/internal/handler/http"
	"github.com/arinamklvch/go-courses/internal/pb"
	"github.com/arinamklvch/go-courses/internal/usecase"
	"google.golang.org/grpc"
)

// @title Movies Downloader API
// @version 1.0
// @description HTTP API for loading movie ranges and finding the last available movie.
func main() {
	err := run()
	if err != nil {
		log.Printf("server stopped with error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.LoadConfig()
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	movieClient := adapter.New(httpClient, cfg.MovieBaseURL)
	movieUseCase := usecase.New(movieClient)
	handler := httphandler.New(movieUseCase)
	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: handler.Routes(),
	}
	defer func() {
		_ = srv.Close()
	}()

	grpcListener, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}
	defer func() {
		_ = grpcListener.Close()
	}()

	grpcServer := grpc.NewServer()
	pb.RegisterMoviesServiceServer(grpcServer, grpchandler.New(movieUseCase))
	defer grpcServer.Stop()

	serverErrChan := make(chan error, 2)
	go func() {
		log.Printf("HTTP server starting on %s", srv.Addr)
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrChan <- err
		}
	}()
	go func() {
		log.Printf("gRPC server started on %s", grpcListener.Addr())
		err := grpcServer.Serve(grpcListener)
		if err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			serverErrChan <- err
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signalChan)

	select {
	case err := <-serverErrChan:
		return err

	case sig := <-signalChan:
		log.Printf("received signal %s; stopping HTTP and gRPC servers", sig)
		grpcServer.Stop()
		log.Printf("gRPC server stopped")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			return err
		}
		log.Printf("HTTP server stopped")
	}

	return nil
}
