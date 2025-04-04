// elastic/internal/repositories/data_loader.go
package repositories

import (
	"context"
	"fmt"
	"strconv"

	"github.com/lonmouth/elastic/internal/data"
	"github.com/olivere/elastic/v7"
)

// LoadDataToElasticsearch загружает данные с использованием Bulk API
func LoadDataToElasticsearch(client *elastic.Client, records []data.Record) error {
	bulkRequest := client.Bulk().Index("places")
	for _, record := range records {
		req := elastic.NewBulkIndexRequest().
			Id(strconv.Itoa(record.ID)).
			Doc(record)
		bulkRequest.Add(req)
	}

	_, err := bulkRequest.Do(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("Данные успешно загружены в Elasticsearch")
	return nil
}

// LoadDataIfIndexNotExists загружает данные только если индекс не существует
func LoadDataIfIndexNotExists(client *elastic.Client, csvFilePath string, schemaFilePath string) error {
	exists, err := client.IndexExists("places").Do(context.Background())
	if err != nil {
		return err
	}

	if exists {
		fmt.Println("Индекс уже существует, загрузка данных пропущена")
		return nil
	}

	csvRecords, err := data.ReadCSVFromFile(csvFilePath)
	if err != nil {
		return err
	}

	records, err := data.ConvertCSVToRecords(csvRecords)
	if err != nil {
		return err
	}

	schema, err := data.ReadSchemaFromFile(schemaFilePath)
	if err != nil {
		return err
	}

	err = CreateIndexAndMapping(client, "places", schema)
	if err != nil {
		return err
	}

	return LoadDataToElasticsearch(client, records)
}
