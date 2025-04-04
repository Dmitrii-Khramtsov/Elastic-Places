// elastic/internal/repositories/elasticsearch.go
package repositories

import (
	"context"
	"encoding/json"

	"github.com/lonmouth/elastic/internal/entities"
	"github.com/olivere/elastic/v7"
)

// PlacesRepository интерфейс для взаимодействия с хранилищем мест
type PlacesRepository interface {
	GetPlaces(limit, offset int) ([]entities.Place, int, error)
	GetClosestPlaces(lat, lon float64, limit int) ([]entities.Place, error)
}

type ElasticsearchStore struct {
	client *elastic.Client
	index  string
}

// NewElasticsearchStore создаёт новый экземпляр ElasticsearchStore
func NewElasticsearchStore(client *elastic.Client, index string) *ElasticsearchStore {
	return &ElasticsearchStore{client: client, index: index}
}

// GetPlaces получает список мест из Elasticsearch с пагинацией
func (s *ElasticsearchStore) GetPlaces(limit, offset int) ([]entities.Place, int, error) {
	// создаем запрос, который возвращает все документы
	query := elastic.NewMatchAllQuery()

	// выполняем поиск в Elasticsearch
	searchResult, err := s.client.Search().
		Index(s.index).          // указываем индекс для поиска
		Query(query).            // указываем запрос для поиска
		From(offset).            // указываем смещение для пагинации
		Size(limit).             // указываем количество записей для пагинации
		Pretty(true).            // форматирование JSON из строчного в читабельный вид с отступами и переносами
		TrackTotalHits(true).    // включаем подсчет всех попаданий
		Do(context.Background()) // выполняем запрос

	if err != nil {
		return nil, 0, err
	}

	// total - общее колличество совпадений
	total := int(searchResult.TotalHits())

	// создаем срез для хранения мест
	var places []entities.Place
	// проходим по всем найденным документам
	for _, hit := range searchResult.Hits.Hits {
		var place entities.Place
		// декодируем JSON в структуру Place
		err := json.Unmarshal(hit.Source, &place)
		if err != nil {
			return nil, 0, err
		}
		// добавляем место в срез
		places = append(places, place)
	}

	// возвращаем список мест, общее количество документов и nil в качестве ошибки
	return places, total, nil
}

// GetClosestPlaces получает список ближайших мест из Elasticsearch
func (s *ElasticsearchStore) GetClosestPlaces(lat, lon float64, limit int) ([]entities.Place, error) {
	// создаём запрос в Elasticsearch, который возвращает все документы
	query := elastic.NewMatchAllQuery()

	// настраиваем сортировку по географическому расстоянию
	sort := elastic.NewGeoDistanceSort("location").
		Point(lat, lon).     // это широта и долгота, относительно которых будет выполняться сортировка
		Order(true).         // устанавливает порядок сортировки. true означает сортировку по возрастанию (ближайшие места будут первыми), а false — по убыванию (дальние места будут первыми)
		Unit("km").          // устанавливает единицы измерения расстояния. В данном случае используются километры (km). Другие возможные значения включают мили (mi), метры (m)
		DistanceType("arc"). // устанавливает тип расстояния. arc означает, что расстояние будет измеряться по дуге большого круга (что является наиболее точным способом измерения расстояния между двумя точками на поверхности Земли). Другие возможные значения включают plane (прямолинейное расстояние) и sloppy_arc (менее точное измерение по дуге)
		IgnoreUnmapped(true) // указывает, что документы, у которых нет поля location, должны быть проигнорированы при сортировке. Это полезно, если в индексе могут быть документы без географических координат

	// Выполняем поиск в Elasticsearch
	searchResult, err := s.client.Search().
		Index(s.index).          // указываем индекс для поиска
		Query(query).            // указываем запрос для поиска
		SortBy(sort).            // указываем сортировку
		Size(limit).             // указываем колличество записей для пагинации
		Pretty(true).            // форматирование JSON из строчного в читабельный вид с отступами и переносами
		Do(context.Background()) // выполняем запрос

	if err != nil {
		return nil, err
	}

	// создаём срез для хранения мест
	var places []entities.Place
	// проходим по всем найденным документам
	for _, hit := range searchResult.Hits.Hits {
		var place entities.Place
		// декодируем JSON в структуру Place
		err := json.Unmarshal(hit.Source, &place)
		if err != nil {
			return nil, err
		}
		// добавляем место в срез
		places = append(places, place)
	}

	// возвращаем список мест и nil в качестве ошибки
	return places, nil
}
