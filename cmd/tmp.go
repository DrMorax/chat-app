package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Message struct {
	Text string `json:"text"`
}

func main() {
	app := fiber.New()

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		log.Println(c.Locals("allowed"))  // true
		log.Println(c.Params("id"))       // 123
		log.Println(c.Query("v"))         // 1.0
		log.Println(c.Cookies("session")) // ""

		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)

			response := Message{Text: fmt.Sprintf("I received your message %s", msg)}
			jsonResponse, err := json.Marshal(response)
			if err != nil {
				log.Println("marshal:", err)
				break
			}

			if err = c.WriteMessage(mt, jsonResponse); err != nil {
				log.Println("write:", err)
				break
			}
		}
	}))

	log.Fatal(app.Listen(":3001"))
}
