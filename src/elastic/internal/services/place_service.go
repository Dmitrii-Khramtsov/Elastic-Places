// elastic/internal/services/place_service.go
package services

import (
	"github.com/lonmouth/elastic/internal/entities"
	"github.com/lonmouth/elastic/internal/repositories"
)

// PlaceService структура, представляющая сервис для работы с местами
type PlaceService struct {
	store repositories.PlacesRepository
}

// NewPlaceService создает новый экземпляр PlaceService
func NewPlaceService(store repositories.PlacesRepository) *PlaceService {
	return &PlaceService{store: store}
}

// GetPlaces получает список мест для указанной страницы
func (s *PlaceService) GetPlaces(limit, offset int) ([]entities.Place, int, error) {
	return s.store.GetPlaces(limit, offset)
}

// GetClosestPlaces получает список ближайших мест
func (s *PlaceService) GetClosestPlaces(lat, lon float64, limit int) ([]entities.Place, error) {
	return s.store.GetClosestPlaces(lat, lon, limit)
}
