package event

import (
	"tn/internal/domain/models"
	"tn/internal/storage"
)

func GetEventCategoryOmitFields(event models.EventCategory) []string {
	omits := make([]string, 0)
	omits = append(omits, "id", "event_id")
	if event.Category == storage.DefaultEmptyStr {
		omits = append(omits, "category")
	}
	if event.Price == storage.DefaultEmptyInt {
		omits = append(omits, "price")
	}
	return omits
}

func GetEventOmitFields(event models.Event) []string {
	omits := make([]string, 0)
	omits = append(omits, "id", "owner_id", "categories")

	if event.Name == storage.DefaultEmptyStr {
		omits = append(omits, "name")
	}

	if event.Description == storage.DefaultEmptyStr {
		omits = append(omits, "description")
	}

	if event.Country == storage.DefaultEmptyStr {
		omits = append(omits, "country")
	}

	if event.City == storage.DefaultEmptyStr {
		omits = append(omits, "city")
	}

	if event.Place == storage.DefaultEmptyStr {
		omits = append(omits, "place")
	}

	if event.Date.IsZero() {
		omits = append(omits, "date")
	}
	if event.Age == storage.DefaultEmptyStr {
		omits = append(omits, "age")
	}

	return omits
}
