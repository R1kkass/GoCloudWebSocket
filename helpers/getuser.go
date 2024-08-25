package helpers

import (
	"errors"
	"mypackages/db"
	Model "mypackages/models"
	"strings"

	"golang.org/x/net/websocket"
)

func GetUser(ws *websocket.Conn) (*Model.User, error) {
	token := ws.Request().Header["Access-Token"]
	token = strings.Split(token[0], " ")

	email := ParseJWT(token[1])

	var user Model.User

	r := db.DB.Model(&Model.User{}).Where("email = ?", email).First(&user)

	if r.RowsAffected == 0 {
		return nil, errors.New("пользователь не найден")
	}

	return &user, nil
}
