// elastic/internal/transport/handlers/place_handler_auth.go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// секретный ключ для подписи JWT токенов
var jwtKey = []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNjAxOTc1ODI5LCJuYW1lIjoiTmlrb2xheSJ9.FqsRe0t9YhvEC3hK1pCWumGvrJgz9k9WvhJgO8HsIa8")

// структура для ответа с токеном
type TokenResponse struct {
	Token string `json:"token"`
}

// HandleGetToken - обрабатывает запрос на получение JWT токенов
func (h *PlaceHandler) HandleGetToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // устанавливаем заголовок Content-Type

	// создаём Claims (утверждения) для токена
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 1).Unix(), // время истечения токена 1 час
		Issuer:    "Elastic",                            // издатель токена
	}

	// создаём новый токен с методом подписи HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// подписываем токен секретным ключом
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // возвращаем ошибку, если подписание не удалось
		return
	}

	// создаём ответ с токеном
	response := TokenResponse{
		Token: tokenStr,
	}

	// кодируем ответ в JSON и отправляем его
	json.NewEncoder(w).Encode(response)
}

// JWTMiddleware является промежуточным ПО для проверки JWT токена в заголовке Authorization
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// получаем значение заголовка Authorization из запроса
		authHeader := r.Header.Get("Authorization")

		// проверяем, что заголовок Authorization присутствует
		if authHeader == "" {
			// Если заголовок отсутствует, возвращаем ошибку 401 Unauthorized
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// удаляем префикс "Bearer " из заголовка и очищаем пробелы
		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))

		// парсим токен с использованием секретного ключа
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// проверяем, что метод подписи токена является HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				// если метод подписи не HMAC, возвращаем ошибку
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// возвращаем секретный ключ для проверки подписи токена
			return jwtKey, nil
		})

		// проверяем, что токен был успешно распарсен
		if err != nil {
			// если произошла ошибка при парсинге токена, возвращаем ошибку 401 Unauthorized
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// проверяем, что токен действителен
		if !token.Valid {
			// если токен недействителен, возвращаем ошибку 401 Unauthorized
			http.Error(w, "Token is not valid", http.StatusUnauthorized)
			return
		}

		// если токен действителен, передаем управление следующему обработчику
		next.ServeHTTP(w, r)
	})
}
