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
	"github.com/arinamklvch/go-courses/internal/config"
	grpccontroller "github.com/arinamklvch/go-courses/internal/controller/grpc"
	httpcontroller "github.com/arinamklvch/go-courses/internal/controller/http"
	"github.com/arinamklvch/go-courses/internal/pb"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"google.golang.org/grpc"
)

func run() error {
	cfg := config.LoadConfig()
	mux := http.NewServeMux()
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	s := httpcontroller.New(client, cfg.MovieBaseURL)

	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	mux.HandleFunc("GET /load-movies", s.GetMoviesRange)
	mux.HandleFunc("GET /last-movie", s.GetLastMovie)

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: mux,
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
	pb.RegisterMoviesServiceServer(grpcServer, grpccontroller.New(client, cfg.MovieBaseURL))
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
