package server

import (
	"log"
	"net"

	"github.com/th3g00dc0d3r/codebank/infrastructure/grpc/pb"
	"github.com/th3g00dc0d3r/codebank/infrastructure/grpc/service"
	"github.com/th3g00dc0d3r/codebank/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	ProcessTransactionUseCase usecase.UseCaseTransaction
}

func NewGRPCServer() GRPCServer {
	return GRPCServer{}
}

func (g GRPCServer) Serve() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Could not listen to tcp port")
	}

	transactionService := service.NewTransactionService()

	transactionService.ProcessTransactionUseCase = g.ProcessTransactionUseCase

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterPaymentServiceServer(grpcServer, transactionService)
	grpcServer.Serve(lis)
}
