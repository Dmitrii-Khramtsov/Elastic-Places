// elastic/internal/transport/handlers/place_handler_utils.go
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/lonmouth/elastic/internal/entities"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// вспомогательная функция для HandleGetPlacesHTML и HandleGetPlacesJSON
func (h *PlaceHandler) fetchPlaces(w http.ResponseWriter, r *http.Request) ([]entities.Place, int, int, int, error) {
	// получаем общее колличество позиций
	_, total, err := h.placeService.GetPlaces(0, 1)
	if err != nil {
		log.Println("Ошибка при получении данных из Elastisearch")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, 0, 0, 0, err
	}

	// вычисляем последнюю страницу
	lastPage := (total + pageSize - 1) / pageSize

	// получаем номер страницы из параметров запроса
	page, err := pageNumber(w, r, lastPage)
	if err != nil {
		log.Printf("Неверный номер страницы: %v", err)
		return nil, 0, 0, 0, err
	}

	// вычисляем смещение для запроса к базе данных
	offset := (page - 1) * pageSize

	// получаем данные о местах для указанной страницы
	places, _, err := h.placeService.GetPlaces(pageSize, offset)
	if err != nil {
		log.Println("Ошибка при получении данных из Elastisearch")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, 0, 0, 0, err
	}

	return places, total, page, lastPage, nil
}

func pageNumber(w http.ResponseWriter, r *http.Request, lastPage int) (int, error) {
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 || page > lastPage {
		w.WriteHeader(http.StatusBadRequest) // устанавливаем статус код сначала

		errorResponse := ErrorResponse{
			Error: fmt.Sprintf("Invalid 'page' value: %s", pageStr), // сообщение об ошибке для клиента
		}
		json.NewEncoder(w).Encode(errorResponse) // записываем сообщение об ошибке в ответ

		return 0, fmt.Errorf("invalid page number: %s", pageStr) // внутренняя ошибка для обработки
	}

	return page, nil
}

func convertDataToJSON(w http.ResponseWriter, places []entities.Place, total, page, lastPage int) error {
	// подготавливаем данные для JSON
	data := TemplateData{
		Places:   places,
		Total:    total,
		Page:     page,
		LastPage: lastPage,
		PrevPage: page - 1,
		NextPage: page + 1,
	}

	// кодируем данные в JSON и отправляем их в ответ
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return err
}
