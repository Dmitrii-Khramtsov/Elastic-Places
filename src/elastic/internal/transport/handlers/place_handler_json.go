// elastic/internal/transport/handlers/place_handler_json.go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"log"

	"github.com/lonmouth/elastic/internal/entities"
)

type RecomendationResponse struct {
	Name   string          `json:"name"`
	Places []entities.Place `json:"places"`
}

// HandleGetClosestPlaces обрабатывает HTTP-запросы, поступающих на определенный маршрут, с целью получения ближайших мест относительно заданных координат
func (h *PlaceHandler) HandleGetClosestPlaces(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// получаем координаты из параметров запроса
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		http.Error(w, "Неверный формат широты", http.StatusBadRequest)
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		http.Error(w, "Неверный формат долготы", http.StatusBadRequest)
	}

	// получаем ближайшие места
	places, err := h.placeService.GetClosestPlaces(lat, lon, limitLocations)
	if err != nil {
		log.Println("Ошибка при получении данных из Elasticsearch")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// создаем структуру данных для передачи в JSON
	response := RecomendationResponse{
		Name:   "Recomendation",
		Places: places,
	}

	// кодируем данные в JSON и передаём в ответ
	err = json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("ответ отправлен успешно")
}

// HandleGetPlacesJSON обрабатывает HTTP-запросы для получения и отображения списка мест в формате JSON
func (h *PlaceHandler) HandleGetPlacesJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	places, total, page, lastPage, err := h.fetchPlaces(w, r)
	if err != nil {
		return
	}

	err = convertDataToJSON(w, places, total, page, lastPage)
	if err != nil {
		log.Println("Ошибка при кодировании данных в JSON")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Ответ отправлен успешно")
}
