package cache

import (
	"errors"
	"event/pkg/models"
	"sync"
	"time"
)

type EventsCache struct {
	events map[uint64]*models.Event
}

func NewEventsCache(events map[uint64]*models.Event) *EventsCache {
	return &EventsCache{events: events}
}

func (e *EventsCache) GetEvent(id uint64) (*models.Event, error) {
	var (
		m     sync.Mutex
		event *models.Event
		ok    bool
	)
	m.Lock()
	if event, ok = e.events[id]; !ok {
		return nil, errors.New("no such event")
	}
	m.Unlock()
	return event, nil
}
func (e *EventsCache) getEventsForN(n float64) ([]*models.Event, error) {
	var (
		currentTime = time.Now()
		events      = make([]*models.Event, 0)
	)
	for _, event := range e.events {
		difference := event.Date.Sub(currentTime).Hours() / 24 / n
		if event.Date.After(currentTime) && difference < 1 {
			events = append(events, event)
		}
	}
	return events, nil
}
func (e *EventsCache) GetEventsForDay() ([]*models.Event, error) {
	return e.getEventsForN(1)
}
func (e *EventsCache) GetEventsForWeek() ([]*models.Event, error) {
	return e.getEventsForN(7)
}
func (e *EventsCache) GetEventsForMonth() ([]*models.Event, error) {
	return e.getEventsForN(30)
}
func (e *EventsCache) AddEvent(event *models.Event) error {
	var (
		m sync.Mutex
	)
	m.Lock()
	e.events[event.Id] = event
	m.Unlock()
	return nil
}
func (e *EventsCache) UpdateEvent(event *models.Event) error {
	var (
		m  sync.Mutex
		ok bool
	)
	m.Lock()
	if event, ok = e.events[event.Id]; !ok {
		return errors.New("no such event")
	}
	e.events[event.Id] = event
	m.Unlock()
	return nil
}
func (e *EventsCache) DeleteEvent(id uint64) error {
	var (
		m  sync.Mutex
		ok bool
	)
	m.Lock()
	if _, ok = e.events[id]; !ok {
		return errors.New("no such event")
	}
	delete(e.events, id)
	m.Unlock()
	return nil
}
