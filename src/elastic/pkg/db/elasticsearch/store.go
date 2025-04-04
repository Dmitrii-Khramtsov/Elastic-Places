// elastic/pkg/db/elasticsearch/store.go
package elasticsearch

import (
	"log"

	"github.com/lonmouth/elastic/internal/repositories"
	"github.com/olivere/elastic/v7"
)

const elasticURL = "http://localhost:9200"

// NewStore создаёт новый экземпляр хранилища
func NewStore() repositories.PlacesRepository {
	client, err := elastic.NewClient(elastic.SetURL(elasticURL))
	if err != nil {
		log.Fatalf("ошибка создания клиента Elasticsearch: %s", err)
	}

	// Загрузка данных только если индекс не существует
	err = repositories.LoadDataIfIndexNotExists(client, "internal/data/data.csv", "internal/schema/schema.json")
	if err != nil {
		log.Fatalf("Ошибка при загрузке данных: %v\n", err)
	}

	return repositories.NewElasticsearchStore(client, "places")
}
