package mywebsockets

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mypackages/db"
	"mypackages/helpers"
	Model "mypackages/models"
	"strconv"

	"golang.org/x/net/websocket"
)

func (s *Server) SendMessages(ws *websocket.Conn, user *Model.User) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read error:", err)
			continue
		}
		msg := buf[:n]

		chatId := ws.Request().URL.Query()["id"][0] 

		chatInt, err := strconv.Atoi(chatId)

		if err != nil {
			ws.WriteClose(400)
		}

		message := &Model.Message{
			Text: string(msg),
			UserRelation: Model.UserRelation{
				UserID: int(user.ID),
			},
			ChatID: chatInt,
		}

		r := db.DB.Preload("User").Create(&message).First(&message)
		
		mapMessage := map[string]any{
			"description": string(msg),
			"title": "Новое сообщение",
			"type": "New_Message",
			"options": map[string]string{
				"chat_id": chatId,
			},
		}
		objectMessage, _ := json.Marshal(mapMessage)
		
		if r.RowsAffected == 0 {
			ws.Write([]byte("Ошибка"))
		} else {
			go unReadedMessages(uint(chatInt), message.ID, int(user.ID))
			jsonMessage, _ := json.Marshal(message)
			go s.sendNotification(chatInt, user.ID, objectMessage)
			s.broadcast(jsonMessage, chatId)
		}
	}
}

func MiddlewareMessage(ws *websocket.Conn) (*Model.User, error) {
	if err := helpers.VerifyJWT(ws); err != nil {
		return nil, err
	}

	user, err := helpers.GetUser(ws)

	if err != nil {
		return nil, err
	}

	if err := CheckChat(ws.Request().URL.Query()["id"][0], user.ID); err != nil {
		return nil, err
	}	
	
	return user, nil
}

func (s *Server) HandleChatWs(ws *websocket.Conn) {

	if err := helpers.VerifyJWT(ws); err != nil {
		ws.WriteClose(401)
	}

	user, err := MiddlewareMessage(ws)
	if err != nil {
		ws.WriteClose(404)
	}

	var messages []*Model.Message

	db.DB.Model(&Model.Message{}).Preload("User").Where("chat_id=?", ws.Request().URL.Query()["id"][0]).Order("id DESC").Limit(10).Find(&messages)
	messageJson, err := json.Marshal(messages)

	if err != nil {
		ws.WriteClose(500)
	}

	ws.Write([]byte(string(messageJson)))
	s.conns[ws] = true

	s.SendMessages(ws, user)
}

func CheckChat(idChat string, userId uint) error {
	var chat []*Model.ChatUser

	result := db.DB.Model(&Model.ChatUser{}).Where("chat_id = ? AND user_id = ?", idChat, userId).Find(&chat)

	if result.RowsAffected == 0 {
		return errors.New("чат не найден")
	}

	for _, item := range chat {
		if !item.SubmitCreate {
			return errors.New("пользователи не подтверждены")
		}
	}

	return nil
}