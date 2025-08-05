package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"dual-server-example/server"
)

func main() {
	grpcPort := ":9000"
	httpPort := ":8080"

	// gRPCサーバーを起動
	go func() {
		if err := server.StartGRPCServer(grpcPort); err != nil {
			log.Fatalf("Failed to start gRPC server: %v", err)
		}
	}()

	// HTTPサーバーを起動
	httpServer, err := server.NewHTTPServer("localhost" + grpcPort)
	if err != nil {
		log.Fatalf("Failed to create HTTP server: %v", err)
	}

	go func() {
		if err := httpServer.StartHTTPServer(httpPort); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	fmt.Printf("🚀 Dual server started!\n")
	fmt.Printf("📡 gRPC server: localhost%s\n", grpcPort)
	fmt.Printf("🌐 HTTP server: localhost%s\n", httpPort)

	// グレースフルシャットダウン
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n🛑 Shutting down servers...")
}
