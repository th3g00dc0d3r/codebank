package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/th3g00dc0d3r/codebank/infrastructure/grpc/server"
	"github.com/th3g00dc0d3r/codebank/infrastructure/kafka"
	"github.com/th3g00dc0d3r/codebank/infrastructure/repository"
	"github.com/th3g00dc0d3r/codebank/usecase"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	db := setupDb()
	defer db.Close()

	producer := setupKafkaProducer()
	
	processTransactionUseCase := setupTransactionUseCase(db, producer)

	serveGrpc(processTransactionUseCase)

}

func setupDb() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("host"),
		os.Getenv("port"),
		os.Getenv("user"),
		os.Getenv("password"),
		os.Getenv("dbname"),
	)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal("ERROR CONNECTING TO DATABASE")
	}

	return db
}

func setupKafkaProducer() kafka.KafkaProducer {
	producer := kafka.NewKafkaProducer()
	producer.SetupProducer(os.Getenv("KafkaBootstrapServers"))

	return producer
}

func setupTransactionUseCase(db *sql.DB, producer kafka.KafkaProducer) usecase.UseCaseTransaction {
	transactionRepository := repository.NewTransactionRepositoryDb(db)
	useCase := usecase.NewUseCaseTransaction(transactionRepository)
	useCase.KafkaProducer = producer

	return useCase
}

func serveGrpc(processTransactionUseCase usecase.UseCaseTransaction) {
	grpcServer := server.NewGRPCServer()

	grpcServer.ProcessTransactionUseCase = processTransactionUseCase

	grpcServer.Serve()
}