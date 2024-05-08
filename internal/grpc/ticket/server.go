package ticketgrpc

import (
	"context"
	"tn/internal/domain/models"

	ticketv1 "github.com/pvdiploma/diploma-protos/gen/go/ticket"
	"google.golang.org/grpc"
)

type Ticket interface {
	AddTicket(ctx context.Context, event_category_id int64, name string, surname string, patronymic string, email string) (int64, error)
	GetTicket(ctx context.Context, ticketID int64) (models.Ticket, error)
	DeleteTicket(ctx context.Context, ticketID int64) (int64, error)
	ActivateTicket(ctx context.Context, ticketID int64) (int64, error)
	IsActivated(ctx context.Context, ticketID int64) (bool, error)
}

type serverAPI struct {
	ticketv1.UnimplementedTicketServiceServer
	ticket Ticket
}

func Register(gRPC *grpc.Server, ticket Ticket) {
	ticketv1.RegisterTicketServiceServer(gRPC, &serverAPI{ticket: ticket})
}

func (s *serverAPI) AddTicket(ctx context.Context, req *ticketv1.AddTicketRequest) (*ticketv1.AddTicketResponse, error) {
	// проверка безопасности ?

	// сразу вызов в сервисный слой с данными вида models.Ticket

	return &ticketv1.AddTicketResponse{}, nil
}

func (s *serverAPI) GetTicketImage(ctx context.Context, req *ticketv1.GetTicketImageRequest) (*ticketv1.GetTicketImageResponse, error) {

	return &ticketv1.GetTicketImageResponse{}, nil
}

func (s *serverAPI) GetTicketInfo(ctx context.Context, req *ticketv1.GetTicketInfoRequest) (*ticketv1.GetTicketInfoResponse, error) {

	return &ticketv1.GetTicketInfoResponse{}, nil
}

func (s *serverAPI) DeleteTicket(ctx context.Context, req *ticketv1.DeleteTicketRequest) (*ticketv1.DeleteTicketResponse, error) {
	return &ticketv1.DeleteTicketResponse{}, nil

}

func (s *serverAPI) ActivateTicket(ctx context.Context, req *ticketv1.ActivateTicketRequest) (*ticketv1.ActivateTicketResponse, error) {
	return &ticketv1.ActivateTicketResponse{}, nil

}
