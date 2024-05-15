package paymentgrpc

import (
	"context"
	tokenmanager "tn/internal/utils/tokenManager"

	paymentv1 "github.com/pvdiploma/diploma-protos/gen/go/payment"
	"google.golang.org/grpc"
)

type Payment interface {
}

type serverAPI struct {
	paymentv1.UnimplementedPaymentServiceServer
	payment Payment
	tm      *tokenmanager.TokenManager
}

func Register(gRPC *grpc.Server, payment Payment, tm *tokenmanager.TokenManager) {
	paymentv1.RegisterPaymentServiceServer(gRPC, &serverAPI{payment: payment, tm: tm})
}

// данный метод должен участвовать в схеме кафки
func (s *serverAPI) PurchaseTicket(ctx context.Context, req *paymentv1.PurchaseTicketRequest) (*paymentv1.PurchaseTicketResponse, error) {

	// проверить возможность покупки билетОВ

	// создать токен для конкретной покуки

	// Возврат созданного id билета
	return &paymentv1.PurchaseTicketResponse{}, nil
}

// должно идти после создания билета и подтверждения оплаты (пока пропускаем второй шаг)
func (s *serverAPI) SendTicket(ctx context.Context, req *paymentv1.SendTicketRequest) (*paymentv1.SendTicketResponse, error) {
	return &paymentv1.SendTicketResponse{}, nil
}
