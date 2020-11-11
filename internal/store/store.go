package store

import (
	"time"

	"city-hall-lights/internal/model"
)

type Store interface {
	Create(event model.Event) error
	Update(event model.Event) error
	Read(date time.Time) (model.Event, error)
	List(date time.Time) ([]model.Event, error)
	Delete(event model.Event) error
}
