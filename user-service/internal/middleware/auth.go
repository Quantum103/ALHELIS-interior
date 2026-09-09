package middleware

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "userID"

func JWTAuthMiddleware(jwtSecret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				http.Error(
					w,
					`{"error":"Отсутствует заголовок Authorization"}`,
					http.StatusUnauthorized,
				)
				return
			}

			headerParts := strings.SplitN(authHeader, " ", 2)

			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				http.Error(
					w,
					`{"error":"Неверный формат заголовка Authorization"}`,
					http.StatusUnauthorized,
				)
				return
			}

			tokenString := headerParts[1]

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("неподдерживаемый метод подписи")
				}

				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				http.Error(
					w,
					`{"error":"Недействительный или просроченный токен"}`,
					http.StatusUnauthorized,
				)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(
					w,
					`{"error":"Ошибка чтения данных токена"}`,
					http.StatusUnauthorized,
				)
				return
			}

			sub, ok := claims["sub"]
			if !ok {
				http.Error(
					w,
					`{"error":"В токене отсутствует user_id"}`,
					http.StatusUnauthorized,
				)
				return
			}

			var userID int64

			switch value := sub.(type) {
			case float64:
				userID = int64(value)

			case string:
				parsedID, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					http.Error(
						w,
						`{"error":"Некорректный user_id"}`,
						http.StatusUnauthorized,
					)
					return
				}

				userID = parsedID

			default:
				http.Error(
					w,
					`{"error":"Некорректный тип user_id"}`,
					http.StatusUnauthorized,
				)
				return
			}

			ctx := context.WithValue(
				r.Context(),
				UserIDKey,
				userID,
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}
