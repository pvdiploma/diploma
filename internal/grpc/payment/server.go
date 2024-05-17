package paymentgrpc

import (
	"context"
	"tn/internal/domain/models"
	"tn/internal/utils/converter"
	tokenmanager "tn/internal/utils/tokenManager"

	paymentv1 "github.com/pvdiploma/diploma-protos/gen/go/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Payment interface {
	CheckAbilityToBuy(ctx context.Context, purchaseTickets []models.PurchaseTickets) (bool, error)
	CreateTickets(ctx context.Context, purchaseTickets models.PurchaseInfo, purchaseToken string) ([]int64, error)
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
	purchaseInfo := models.PurchaseInfo{
		WidgetID:        req.GetWidgetId(),
		Name:            req.GetName(),
		Surname:         req.GetSurname(),
		Patronymic:      req.GetPatronymic(),
		Email:           req.GetEmail(),
		Phone:           req.GetPhone(),
		PurchaseTickets: converter.ProtoPurchaseTicketsToModels(req.GetCategories()),
	}

	purchaseToken, err := s.tm.NewPurchaseToken()

	// do i need return this message to client?
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create purchase token")
	}

	fl, err := s.payment.CheckAbilityToBuy(ctx, purchaseInfo.PurchaseTickets)

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to check ability to buy")
	}

	if !fl {
		return nil, status.Error(codes.FailedPrecondition, "not enough tickets")
	}

	tikcetsID, err := s.payment.CreateTickets(ctx, purchaseInfo, purchaseToken)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create tickets")
	}

	return &paymentv1.PurchaseTicketResponse{
		Id: tikcetsID,
	}, nil
}

// должно идти после создания билета и подтверждения оплаты (пока пропускаем второй шаг)
func (s *serverAPI) SendTicket(ctx context.Context, req *paymentv1.SendTicketRequest) (*paymentv1.SendTicketResponse, error) {
	return &paymentv1.SendTicketResponse{}, nil
}
