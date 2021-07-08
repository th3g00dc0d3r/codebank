package service

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/th3g00dc0d3r/codebank/dto"
	"github.com/th3g00dc0d3r/codebank/infrastructure/grpc/pb"
	"github.com/th3g00dc0d3r/codebank/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TransactionService struct {
	ProcessTransactionUseCase usecase.UseCaseTransaction
	pb.UnimplementedPaymentServiceServer
}

func NewTransactionService() *TransactionService {
	return &TransactionService{}
}

func (t *TransactionService) Payment(ctx context.Context, in *pb.PaymentRequest) (*empty.Empty, error) {
	transactionDto := dto.Transaction{
		Name:            in.GetCreditcard().GetName(),
		Number:          in.GetCreditcard().GetNumber(),
		ExpirationMonth: in.GetCreditcard().GetExpirationMonth(),
		ExpirationYear:  in.GetCreditcard().GetExpirationYear(),
		CVV:             in.GetCreditcard().GetCVV(),
		Amount:          in.GetAmount(),
		Store:           in.GetStore(),
		Description:     in.GetDescription(),
	}

	transaction, err := t.ProcessTransactionUseCase.ProcessTransaction(transactionDto)

	if err != nil {
		return &empty.Empty{}, status.Error(codes.FailedPrecondition, err.Error())
	}
	if transaction.Status != "approved" {
		return &empty.Empty{}, status.Error(codes.FailedPrecondition, "transaction rejected by the bank")
	}

	return &empty.Empty{}, nil
}
