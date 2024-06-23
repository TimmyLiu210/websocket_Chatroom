package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

type Message struct {
	Event   string `json:"event"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

func NewMessage(event, name, content string) *Message {
	return &Message{
		Event:   event,
		Name:    name,
		Content: content,
	}
}

func (m *Message) GetByteMessage() []byte {
	result, _ := json.Marshal(m)
	return result
}

func main() {

	file, err := os.Create("history.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	r := gin.Default()

	r.LoadHTMLGlob("template/html/*")
	r.Static("/assets", "./template/assets")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	m := melody.New()

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	r.GET("/history", func(c *gin.Context) {
		data, err := ioutil.ReadFile("history.txt")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.String(http.StatusOK, string(data))
	})

	m.HandleConnect(func(s *melody.Session) {
		id := s.Request.URL.Query().Get("id")
		err := writeTxt(file, id, "加入聊天室")
		if err != nil {
			log.Println(err)
		}

		msg := NewMessage("other", id, "加入聊天室").GetByteMessage()
		m.Broadcast(msg)
	})

	m.HandleDisconnect(func(s *melody.Session) {
		id := s.Request.URL.Query().Get("id")
		err := writeTxt(file, id, "離開聊天室")
		if err != nil {
			log.Println(err)
		}

		msg := NewMessage("other", id, "離開聊天室").GetByteMessage()
		m.Broadcast(msg)
	})

	m.HandleMessage(func(s *melody.Session, jsonMsg []byte) {
		id := s.Request.URL.Query().Get("id")

		var msg Message
		err := json.Unmarshal(jsonMsg, &msg)
		if err != nil {
			log.Println(err)
		} else {
			err = writeTxt(file, id, msg.Content)
			if err != nil {
				log.Println(err)
			}
		}

		m.Broadcast(jsonMsg)
	})

	m.HandleError(func(s *melody.Session, err error) {
		log.Print("[error]", err)
	})

	r.Run(":5000")
}

func getTime() string {
	now := time.Now()
	formattedTime := now.Format("[2006:01:02:15:04:05]")
	return formattedTime
}

func writeTxt(file *os.File, id string, context string) error {
	formattedTime := getTime()
	_, err := file.WriteString(formattedTime + id + ":" + context + "\n")
	if err != nil {
		return nil
	}

	return nil
}
