// elastic/internal/transport/handlers/place_handler_html.go
package handlers

import (
	"html/template"
	"net/http"
	"log"

	"github.com/lonmouth/elastic/internal/entities"
)

type TemplateData struct {
	Places   []entities.Place `json:"Places"`
	Total    int              `json:"Total"`
	Page     int              `json:"Page"`
	LastPage int              `json:"LastPage"`
	PrevPage int              `json:"PrevPage"`
	NextPage int              `json:"NextPage"`
}

// HandleGetPlacesHTML обрабатывает HTTP-запросы для получения и отображения списка мест в формате HTML
func (h *PlaceHandler) HandleGetPlacesHTML(w http.ResponseWriter, r *http.Request) {
	places, total, page, lastPage, err := h.fetchPlaces(w, r)
	if err != nil {
		return
	}

	convertDataToHTML(w, places, total, page, lastPage)
	log.Println("Ответ отправлен успешно")
}

func convertDataToHTML(w http.ResponseWriter, places []entities.Place, total, page, lastPage int) {
	// создаем данные для шаблона
	data := struct {
		Places    []entities.Place
		Total     int
		Page      int
		LastPage  int
		ShowFirst bool
		ShowPrev  bool
		ShowNext  bool
	}{
		Places:    places,
		Total:     total,
		Page:      page,
		LastPage:  lastPage,
		ShowFirst: page > 1,
		ShowPrev:  page > 1,
		ShowNext:  page < lastPage,
	}

	// парсим и выполняем шаблон
	tmpl := template.Must(template.New("places").Funcs(template.FuncMap{
		"sub": func(a, b int) int { return a - b },
		"add": func(a, b int) int { return a + b },
	}).Parse(placesTemplate))
	tmpl.Execute(w, data)
}

// шаблон HTML для отображения мест
const placesTemplate = `
<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <title>Места</title>
    <meta name="description" content="">
    <meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
<h5>Всего: {{.Total}}</h5>
<ul>
    {{range .Places}}
    <li>
        <div>{{.Name}}</div>
        <div>{{.Address}}</div>
        <div>{{.Phone}}</div>
    </li>
    {{end}}
</ul>
{{if .ShowFirst}}<a href="/?page=1">First</a>{{end}}
{{if .ShowPrev}}<a href="/?page={{sub .Page 1}}">Prev</a>{{end}}
{{if .ShowNext}}<a href="/?page={{add .Page 1}}">Next</a>{{end}}
{{if ne .Page .LastPage}}<a href="/?page={{.LastPage}}">Last</a>{{end}}
</body>
</html>
`
