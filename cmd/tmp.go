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

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, world")
	})

	app.Get("/ws/:id", websocket.New(func(ctx *websocket.Conn) {
		log.Println(ctx.Locals("allowed"))  // true
		log.Println(ctx.Params("id"))       // 123
		log.Println(ctx.Query("v"))         // 1.0
		log.Println(ctx.Cookies("session")) // ""

		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = ctx.ReadMessage(); err != nil {
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

			if err = ctx.WriteMessage(mt, jsonResponse); err != nil {
				log.Println("write:", err)
				break
			}
		}
	}))

	log.Fatal(app.Listen(":3001"))
}
