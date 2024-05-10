package ticketgrpc

import (
	"context"
	"errors"
	"tn/internal/domain/models"
	"tn/internal/storage"

	ticketv1 "github.com/pvdiploma/diploma-protos/gen/go/ticket"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Ticket interface {
	AddTicket(ctx context.Context, event_category_id int64, name string, surname string, patronymic string, discount uint32, email string) (int64, error)
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

	id, err := s.ticket.AddTicket(ctx, req.GetEventCategoryId(), req.GetName(), req.GetSurname(), req.GetPatronymic(), req.GetDiscount(), req.GetEmail())
	if err != nil {
		if errors.Is(err, storage.ErrTicketExists) {
			return nil, status.Error(codes.AlreadyExists, "ticket already exists")
		}

		if errors.Is(err, storage.ErrEventNotFound) {
			return nil, status.Error(codes.NotFound, "event not found")
		}
		return nil, status.Error(codes.Internal, "failed to add ticket")
	}

	return &ticketv1.AddTicketResponse{
		Id: id,
	}, nil
}

func (s *serverAPI) GetTicketImage(ctx context.Context, req *ticketv1.GetTicketImageRequest) (*ticketv1.GetTicketImageResponse, error) {
	ticket, err := s.ticket.GetTicket(ctx, req.GetId())

	if err != nil {
		if errors.Is(err, storage.ErrTicketNotFound) {
			return nil, status.Error(codes.NotFound, "ticket not found")
		}
		return nil, status.Error(codes.Internal, "failed to get ticket")
	}

	return &ticketv1.GetTicketImageResponse{
		Image: ticket.ImageBytes,
	}, nil
}

func (s *serverAPI) GetTicketInfo(ctx context.Context, req *ticketv1.GetTicketInfoRequest) (*ticketv1.GetTicketInfoResponse, error) {
	// not implemented
	return &ticketv1.GetTicketInfoResponse{}, nil
}

func (s *serverAPI) DeleteTicket(ctx context.Context, req *ticketv1.DeleteTicketRequest) (*ticketv1.DeleteTicketResponse, error) {
	id, err := s.ticket.DeleteTicket(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, storage.ErrTicketNotFound) {
			return nil, status.Error(codes.NotFound, "ticket not found")
		}
		return nil, status.Error(codes.Internal, "failed to delete ticket")
	}
	return &ticketv1.DeleteTicketResponse{
		Id: id,
	}, nil
}

func (s *serverAPI) ActivateTicket(ctx context.Context, req *ticketv1.ActivateTicketRequest) (*ticketv1.ActivateTicketResponse, error) {
	//TODO: implement dat shit
	return &ticketv1.ActivateTicketResponse{}, nil
}
