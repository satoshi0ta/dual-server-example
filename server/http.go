package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "dual-server-example/proto"
)

type HTTPServer struct {
	grpcClient pb.BookServiceClient
}

func NewHTTPServer(grpcAddr string) (*HTTPServer, error) {
	conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %v", err)
	}

	client := pb.NewBookServiceClient(conn)
	return &HTTPServer{grpcClient: client}, nil
}

func (s *HTTPServer) StartHTTPServer(port string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /books", s.handleListBooks)
	mux.HandleFunc("GET /books/{bookId}", s.handleGetBook)
	mux.HandleFunc("POST /books", s.handleCreateBook)

	fmt.Printf("HTTP server listening on %s\n", port)
	return http.ListenAndServe(port, mux)
}

func (s *HTTPServer) handleListBooks(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	resp, err := s.grpcClient.ListBooks(ctx, &pb.ListBooksRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *HTTPServer) handleGetBook(w http.ResponseWriter, r *http.Request) {
	bookId := r.PathValue("bookId")
	if bookId == "" {
		http.Error(w, "book ID is required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	resp, err := s.grpcClient.GetBook(ctx, &pb.GetBookRequest{BookId: bookId})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *HTTPServer) handleCreateBook(w http.ResponseWriter, r *http.Request) {
	var book pb.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	resp, err := s.grpcClient.CreateBook(ctx, &pb.CreateBookRequest{Book: &book})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
