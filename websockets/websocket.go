package mywebsockets

import (
	"context"
	"fmt"
	"mypackages/db"
	"mypackages/helpers"
	Model "mypackages/models"
	"strconv"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func (s *Server) broadcast(b []byte, room string) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if room == ws.Request().URL.Query()["id"][0] {
				_, err := ws.Write(b)
				if err != nil {
					delete(s.conns, ws)
					ws.Close()
					fmt.Println("write error:", s.conns[ws])
				}
			}
		}(ws)
	}
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) sendNotification(chatId int, userId uint, objectMessage any) {
	var users []*Model.ChatUser
	var usersObjectID = make(map[int]bool)
	db.DB.Model(&Model.ChatUser{}).Where("chat_id = ?", chatId).Find(&users)
	for ws := range s.conns {
		user, err := helpers.GetUser(ws)
		if err != nil {
			fmt.Println("Can not get user: ", err)
			break
		}
		fmt.Println(user.ID)
		usersObjectID[int(user.ID)] = true
	}
	for _, user := range users {
		_, ok := usersObjectID[user.UserID]
		if user.ID != userId && !ok {
			db.ConnectRedisDB.Publish(context.TODO(), strconv.Itoa(int(user.UserID)) + "_notification", objectMessage)
		}
	}
}

func unReadedMessages(chatId uint, messageId uint, currentUserId int) {
	var chatUsers []*Model.ChatUser

	r := db.DB.Model(&Model.ChatUser{}).Where("chat_id = ?", chatId).Find(&chatUsers)

	if r.RowsAffected == 0 || r.Error != nil {
		return
	}

	for _, user := range chatUsers {
		if currentUserId != user.UserID {
			fmt.Println(user.ID, currentUserId)
			db.DB.Create(&Model.UnReadedMessages{
				ChatID: chatId,
				UserRelation: Model.UserRelation{
					UserID: user.UserID,
				},
				MessageID: messageId,
			})
		}
	}
}