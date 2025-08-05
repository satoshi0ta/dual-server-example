package server

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "dual-server-example/proto"
)

type BookServer struct {
	pb.UnimplementedBookServiceServer
	books map[string]*pb.Book
}

func NewBookServer() *BookServer {
	return &BookServer{
		books: make(map[string]*pb.Book),
	}
}

func (s *BookServer) ListBooks(ctx context.Context, req *pb.ListBooksRequest) (*pb.ListBooksResponse, error) {
	books := make([]*pb.Book, 0, len(s.books))
	for _, book := range s.books {
		books = append(books, book)
	}
	return &pb.ListBooksResponse{Books: books}, nil
}

func (s *BookServer) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.Book, error) {
	book, exists := s.books[req.BookId]
	if !exists {
		return nil, fmt.Errorf("book not found: %s", req.BookId)
	}
	return book, nil
}

func (s *BookServer) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*pb.Book, error) {
	book := req.Book
	if book.BookId == "" {
		book.BookId = fmt.Sprintf("book_%d", len(s.books)+1)
	}
	s.books[book.BookId] = book
	return book, nil
}

func StartGRPCServer(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterBookServiceServer(s, NewBookServer())
	reflection.Register(s)

	fmt.Printf("gRPC server listening on %s\n", port)
	return s.Serve(lis)
}
