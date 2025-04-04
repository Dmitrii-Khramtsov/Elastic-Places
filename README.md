# Elastic-Places — Поиск ближайших заведений

## Оглавление

1. [Описание проекта](#i)
2. [Структура проекта](#ii)
3. [Шаги для запуска](#iii)
4. [Ключевые технологии](#iv)
5. [API Endpoints](#v)
6. [Особенности реализации](#vi)
7. [Примеры запросов](#vii)
8. [Безопасность](#viii)

<h2 id="i" >Описание проекта</h2>
Elastic-Places — сервис для поиска ближайших заведений (ресторанов, кафе) с использованием Elasticsearch.

Особенности:
- Геопоиск по координатам
- Пагинация через 10 результатов
- JWT аутентификация для защищенных запросов
- Поддержка HTML и JSON форматов

<h2 id="ii" >Структура проекта</h2>

```
elastic/
├── cmd/
│   └── elastic/
│       └── main.go
├── internal/
│   ├── data/
│   │   └── data_preparation.go
│   │   └── data.csv
│   ├── entities/
│   │   └── place.go
│   ├── repositories/
│   │   └── data_loader.go
│   │   └── elasticsearch.go
│   │   └── index.go
│   └── schema/
│   │   └── schema.json
│   └── services/
│   │   └── place_service.go
│   └── transport/
│       └── handlers/
│           └── place_handler_auth.go
│           └── place_handler_html.go
│           └── place_handler_json.go
│           └── place_handler_utils.go
│           └── place_handler.go
├── pkg/
│   └── db/
│       └── elasticsearch/
│           └── store.go
├── go.mod
└── go.sum
```

<h2 id="iii" >Шаги для запуска</h2>
1. Требования:
   - Go 1.21+
   - Elasticsearch 7.x
   - Установить зависимости: `go mod tidy`

2. Запуск Elasticsearch:

    ```bash
    ./elasticsearch/bin/elasticsearch
    ```

3. Сборка и запуск:

    ```bash
    go build -o elastic cmd/elastic/main.go
    ./elastic
    ```

4. Веб-интерфейс: http://localhost:8888

<h2 id="iv" >Ключевые технологии</h2>

- Elasticsearch - хранение и геопоиск

- Gorilla Mux - маршрутизация

- JWT - аутентификация

- Bulk API - массовая загрузка данных

- Geo Distance Sort - сортировка по расстоянию

<h2 id="v" >API Endpoints</h2>

  Метод | Путь           | Назначение            | Аутентификация |
 |------|----------------|-----------------------|----------------|
 | GET  | /              | HTML список заведений | Нет            |
 | GET  | /api/places    | JSON пагинация        | Нет            |
 | GET  | /api/recommend | 3 ближайших места     | Да (JWT)       |
 | GET  | /api/get_token | Получить JWT токен    | Нет            |

<h2 id="vi" >Особенности реализации</h2>

- Валидация параметров:

   - lat/lon в диапазоне -90/90 и -180/180

   - page > 0 и <= последней страницы

- Обработка ошибок:

    ``` json
    {"error": "Invalid 'lat' value: 'abc'"}
    ```

- Пагинация через from/size Elasticsearch

- Автоматическое создание индекса при первом запуске

<h2 id="vii" >Примеры запросов</h2>

1. Получение токена:

    ```bash
    curl http://localhost:8888/api/get_token
    ```

    Ответ:

    ```json
    {"token":"eyJhbG...HsIa8"}
    ```

2. Поиск ближайших мест:

    ```bash
    curl -H "Authorization: Bearer <token>" \
    "http://localhost:8888/api/recommend?lat=55.7558&lon=37.6173"
    ```

3. HTML интерфейс: http://localhost:8888/?page=3

<h2 id="viii" >Безопасность</h2>

- JWT секрет хранится в коде (только для демо)

- 401 для неавторизованных запросов к /api/recommend

- 400 для некорректных параметров

- Токен действителен 1 час
