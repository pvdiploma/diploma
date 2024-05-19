package converter

import (
	"tn/internal/domain/models"

	dealv1 "github.com/pvdiploma/diploma-protos/gen/go/deal"
	eventv1 "github.com/pvdiploma/diploma-protos/gen/go/event"
	paymentv1 "github.com/pvdiploma/diploma-protos/gen/go/payment"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ProtoCategoryToModels(reqData []*eventv1.EventCategory) []models.EventCategory {
	var eventCategories []models.EventCategory

	for i := range reqData {
		eventCategories = append(eventCategories, models.EventCategory{
			ID:       reqData[i].GetId(),
			EventID:  reqData[i].GetEventId(),
			Category: reqData[i].GetCategory(),
			Price:    reqData[i].GetPrice(),
			Amount:   reqData[i].GetAmount(),
		})
	}
	return eventCategories
}

func ProtoCategoryToModel(reqData *eventv1.EventCategory) models.EventCategory {
	return models.EventCategory{
		ID:       reqData.GetId(),
		EventID:  reqData.GetEventId(),
		Category: reqData.GetCategory(),
		Price:    reqData.GetPrice(),
		Amount:   reqData.GetAmount(),
	}
}

func ModelCategoryToProto(events []models.EventCategory) []*eventv1.EventCategory {
	var eventCategories []*eventv1.EventCategory

	for i := range events {

		eventCategories = append(eventCategories, &eventv1.EventCategory{
			Id:       events[i].ID,
			EventId:  events[i].EventID,
			Category: events[i].Category,
			Price:    events[i].Price,
			Amount:   events[i].Amount,
		})
	}
	return eventCategories
}

// fix naming ...
func ModelCategoryToProto2(category models.EventCategory) *eventv1.EventCategory {
	return &eventv1.EventCategory{
		Id:       category.ID,
		EventId:  category.EventID,
		Category: category.Category,
		Price:    category.Price,
		Amount:   category.Amount,
	}
}

func ModelEventToProto(event models.Event) *eventv1.Event {
	return &eventv1.Event{
		EventId:      event.ID,
		OwnerId:      event.OwnerID,
		Name:         event.Name,
		Country:      event.Country,
		City:         event.City,
		Place:        event.Place,
		Address:      event.Address,
		Date:         timestamppb.New(event.Date),
		TicketAmount: event.TicketAmount,
		Age:          event.Age,
		Categories:   ModelCategoryToProto(event.Categories),
	}
}

func ProtoEventToModel(event *eventv1.Event) models.Event {
	return models.Event{
		ID:           event.GetEventId(),
		OwnerID:      event.GetOwnerId(),
		Name:         event.GetName(),
		Country:      event.GetCountry(),
		City:         event.GetCity(),
		Place:        event.GetPlace(),
		Date:         event.GetDate().AsTime(),
		TicketAmount: event.GetTicketAmount(),
		Age:          event.GetAge(),
		Categories:   ProtoCategoryToModels(event.GetCategories()),
	}
}

func ProtoPurchaseTicketsToModels(tickets []*paymentv1.TicketCategory) []models.PurchaseTickets {
	var purchaseTickets []models.PurchaseTickets

	for _, ticket := range tickets {
		purchaseTickets = append(purchaseTickets, models.PurchaseTickets{
			ID:     ticket.GetEventCategoryId(),
			Amount: ticket.GetAmount(),
		})
	}
	return purchaseTickets
}

func DealStatusToProto(status models.DealStatus) dealv1.DealStatus {
	switch status {
	case models.Accepted:
		return dealv1.DealStatus_ACCEPTED
	case models.Rejected:
		return dealv1.DealStatus_REJECTED
	default:
		return dealv1.DealStatus_PENDING
	}
}

func ProtoDealStatusToModels(status dealv1.DealStatus) models.DealStatus {
	switch status {
	case dealv1.DealStatus_ACCEPTED:
		return models.Accepted
	case dealv1.DealStatus_REJECTED:
		return models.Rejected
	default:
		return models.Pending
	}
}

func ModelDealsToProto(deals []models.Deal) []*dealv1.Deals {
	var dealsProto []*dealv1.Deals

	for _, deal := range deals {
		dealsProto = append(dealsProto, &dealv1.Deals{
			Id:            deal.ID,
			SenderId:      deal.SenderID,
			RecipientId:   deal.RecipientID,
			OrganizerId:   deal.OrganizerID,
			DistributorId: deal.DistributorID,
			Commission:    deal.Commission,
			EventId:       deal.EventID,
			Status:        DealStatusToProto(deal.Status),
		})
	}
	return dealsProto
}
