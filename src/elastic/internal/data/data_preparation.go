// elastic/internal/data/data_preparation.go
package data

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

// Record структура для хранения данных из CSV
type Record struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Location struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"location"`
}

// ReadCSVFromFile читает CSV файл и заполняет двумерный массив
func ReadCSVFromFile(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

// FilterInvalidRecords фильтрует некорректные записи
func FilterInvalidRecords(records [][]string) ([][]string, error) {
	var validRecords [][]string
	for _, record := range records {
		if len(record) != 6 {
			return nil, fmt.Errorf("invalid record: %v", record)
		}
		validRecords = append(validRecords, record)
	}
	return validRecords, nil
}

// ParseRecords парсит записи
func ParseRecords(records [][]string) ([]Record, error) {
	var data []Record
	for i, record := range records {
		id := i + 1
		lat, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid latitude: %v", record[5])
		}
		lon, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid longitude: %v", record[4])
		}
		data = append(data, Record{
			ID:      id,
			Name:    record[1],
			Address: record[2],
			Phone:   record[3],
			Location: struct {
				Lat float64 `json:"lat"`
				Lon float64 `json:"lon"`
			}{Lat: lat, Lon: lon},
		})
	}
	return data, nil
}

// ConvertCSVToRecords преобразует данные из CSV в структуру
func ConvertCSVToRecords(records [][]string) ([]Record, error) {
	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}
	records = records[1:]

	validRecords, err := FilterInvalidRecords(records)
	if err != nil {
		return nil, err
	}

	data, err := ParseRecords(validRecords)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// ReadSchemaFromFile читает схему маппинга из файла
func ReadSchemaFromFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
