// elastic/internal/repositories/index.go
package repositories

import (
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
)

// CreateIndexAndMapping создает индекс и маппинг
func CreateIndexAndMapping(client *elastic.Client, indexName string, schema string) error {
	exists, err := client.IndexExists(indexName).Do(context.Background())
	if err != nil {
		return err
	}

	if exists {
		fmt.Printf("Индекс %s уже существует\n", indexName)
		return nil
	}

	createIndex, err := client.CreateIndex(indexName).BodyString(schema).Do(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Индекс создан:", createIndex.Acknowledged)

	return nil
}
