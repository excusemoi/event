package repository

import "event/pkg/models"

type EventsDb interface {
	GetEvent(id uint64) (*models.Event, error)
	GetEventsForDay() ([]*models.Event, error)
	GetEventsForWeek() ([]*models.Event, error)
	GetEventsForMonth() ([]*models.Event, error)
	AddEvent(event *models.Event) error
	UpdateEvent(event *models.Event) error
	DeleteEvent(id uint64) error
}
