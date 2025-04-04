// elastic/cmd/elastic/main.go

// elastic/
// ├── cmd/
// │   └── elastic/
// │       └── main.go
// ├── internal/
// │   ├── data/
// │   │   └── data_preparation.go
// │   │   └── data.csv
// │   ├── entities/
// │   │   └── place.go
// │   ├── repositories/
// │   │   └── data_loader.go
// │   │   └── elasticsearch.go
// │   │   └── index.go
// │   └── schema/
// │   │   └── schema.json
// │   └── services/
// │   │   └── place_service.go
// │   └── transport/
// │       └── handlers/
// │           └── place_handler_auth.go
// │           └── place_handler_html.go
// │           └── place_handler_json.go
// │           └── place_handler_utils.go
// │           └── place_handler.go
// ├── pkg/
// │   └── db/
// │       └── elasticsearch/
// │           └── store.go
// ├── go.mod
// └── go.sum

package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lonmouth/elastic/pkg/db/elasticsearch"
	"github.com/lonmouth/elastic/internal/transport/handlers"
	"github.com/lonmouth/elastic/internal/services"
)

func main() {
	store := elasticsearch.NewStore()
	placeService := services.NewPlaceService(store)
	placeHandler := handlers.NewPlaceHandler(placeService)

	router := mux.NewRouter() // создаём новый маршрутизатор
	router.HandleFunc("/", placeHandler.HandleGetPlacesHTML).Methods("GET") // настраиваем маршрутизатор на получение всех элементов в формате HTML
	router.HandleFunc("/api/places", placeHandler.HandleGetPlacesJSON).Methods("GET") // настраиваем маршрутизатор на получение всех элементов в формате JSON

	recommendHandler := http.HandlerFunc(placeHandler.HandleGetClosestPlaces)
	router.Handle("/api/recommend", handlers.JWTMiddleware(recommendHandler)).Methods("GET") // настраиваем маршрутизатор на получение 3-х ближайших мест в формате JSON с аутентификацией

	tokenHandler := http.HandlerFunc(placeHandler.HandleGetToken)
	router.Handle("/api/get_token", tokenHandler).Methods("GET") // настраиваем маршрутизатор на получение токена в формате JSON

	log.Println("сервер запущен на порту 8888")
	log.Fatal(http.ListenAndServe(":8888", router))
}

// curl -X GET "http://127.0.0.1:8888/api/recommend?lat=55.7558&lon=37.6173"
// curl -X GET "http://127.0.0.1:8888/api/recommend?lat=55.7558&lon=37.6173" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzQzOTU4NzcsImlzcyI6IkVsYXN0aWMifQ.q9VWsDNtLUf2lBXoD74_uDCv_dELrZs7UAiA3E1hJ14"

// curl -X DELETE "http://localhost:9200/places"
// curl -XPUT -H "Content-Type: application/json" "http://localhost:9200/places/_settings" -d '
// {
//   "index" : {
//     "max_result_window" : 20000
//   }
// }'
