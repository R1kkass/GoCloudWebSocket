package helpers

import (
	"errors"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/net/websocket"
)

var seccretKey, _ = os.LookupEnv("SECRET_KEY")

func VerifyJWT(r *websocket.Conn) error {
	accessKey := r.Request().Header["Access-Token"]

	accessKey = strings.Split(accessKey[0], " ")

	if len(accessKey) < 2 {
		return errors.New("ошибка токена")
	}

	token, err := jwt.Parse(accessKey[1], func(token *jwt.Token) (interface{}, error) {
		return []byte(seccretKey), nil
	})

	if err != nil {
		return errors.New("ошибка токена")
	}

	if !token.Valid {
		return errors.New("ошибка токена")
	}
	return nil
}
