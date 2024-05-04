package converter

import (
	"tn/internal/domain/models"

	eventv1 "github.com/pvdiploma/diploma-protos/gen/go/event"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ProtoCategoryToModels(reqData []*eventv1.EventCategory) []models.EventCategory {
	var eventCategories []models.EventCategory

	for i := range reqData {
		eventCategories = append(eventCategories, models.EventCategory{
			Category: reqData[i].GetCategory(),
			Price:    reqData[i].GetPrice(),
			Amount:   reqData[i].GetAmount(),
		})
	}
	return eventCategories
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

func ModelEventToProto(event models.Event) *eventv1.Event {
	return &eventv1.Event{
		EventId:      event.ID,
		OwnerId:      event.OwnerID,
		Name:         event.Name,
		Country:      event.Country,
		City:         event.City,
		Place:        event.Place,
		Date:         timestamppb.New(event.Date),
		TicketAmount: event.TicketAmount,
		Age:          event.Age,
		Categories:   ModelCategoryToProto(event.Categories),
	}
}
