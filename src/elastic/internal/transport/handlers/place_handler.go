// elastic/internal/transport/handlers/place_handler.go
package handlers

import (
	"github.com/lonmouth/elastic/internal/services"
	"github.com/lonmouth/elastic/internal/entities"
)

const (
	pageSize       = 10
	limitLocations = 3
)

// PlacesService интерфейс для взаимодействия с сервисом мест
type PlacesService interface {
	GetPlaces(limit, offset int) ([]entities.Place, int, error)
	GetClosestPlaces(lat, lon float64, limit int) ([]entities.Place, error)
}

// PlaceHandler структура для обработки HTTP-запросов
type PlaceHandler struct {
	placeService *services.PlaceService
}

// NewPlaceHandler создает новый экземпляр PlaceHandler
func NewPlaceHandler(placeService *services.PlaceService) *PlaceHandler {
	return &PlaceHandler{placeService: placeService}
}
