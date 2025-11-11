package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	grpcHandler "github.com/lifedaemon-kill/industrial-backend-development-task/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/lifedaemon-kill/industrial-backend-development-task/internal/service"
	desc "github.com/lifedaemon-kill/industrial-backend-development-task/pkg/protogen"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/credentials/insecure"
)

const openApiFilePath = "./api/docs/calc.swagger.json"

func main() {
	fmt.Println("Industrial is starting")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	wg := sync.WaitGroup{}
	wg.Go(func() { startGRPC(ctx) })
	wg.Go(func() { startGateway(ctx) })
	wg.Go(func() { startSwagger(ctx) })

	<-sig
	fmt.Println("SIGINT received, shutting down...")
	cancel()

	wg.Wait()

	fmt.Println("Shutdown complete")
}

func startGRPC(ctx context.Context) {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("failed to listen:", err)
	}
	calcService := service.New()
	grpcH := grpcHandler.New(calcService)
	grpcServer := grpc.NewServer()
	desc.RegisterCalculatorServer(grpcServer, grpcH)

	log.Println("Listening gRPC", "port", 50051)
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Println("failed to serve", "err", err)
			return
		}
	}()

	<-ctx.Done()

	log.Println("shutting down gracefully...")

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		log.Println("gRPC server stopped")
	case <-time.After(5 * time.Second):
		log.Println("timeout, forcing stop")
		grpcServer.Stop()
	}
}

func startGateway(ctx context.Context) {
	gwMux := runtime.NewServeMux()
	err := desc.RegisterCalculatorHandlerFromEndpoint(
		ctx, gwMux, "localhost:50051",
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		log.Println("failed to register gRPC-Gateway handler", "err", err)
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // или конкретно: []string{"http://localhost:8081"}
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	gin.SetMode(gin.ReleaseMode)

	r.Any("/*any", gin.WrapH(gwMux))

	srv := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: r,
	}

	log.Println("Starting HTTP server", "port ", 8080)
	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println("server error", "err", err)
		}
	}()

	<-ctx.Done()

	gsctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = srv.Shutdown(gsctx); err != nil {
		log.Println("server shutdown error", err)
	} else {
		log.Println("shutdown complete")
	}
}

func startSwagger(ctx context.Context) {
	mux := chi.NewMux()

	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		b, err := os.ReadFile(openApiFilePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("failed to read swagger.json:", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})

	mux.HandleFunc("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))

	addr := "0.0.0.0:8081"
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		log.Println("Swagger server listening on", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down Swagger server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Swagger server shutdown failed: %v", err)
	}

	log.Println("Swagger server exited gracefully")
}
