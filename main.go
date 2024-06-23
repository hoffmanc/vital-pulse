package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis/v8"
	fiber "github.com/gofiber/fiber/v2"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	app, _, _ := InitApp()

	port := os.Getenv("PORT")
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}

func InitApp() (*fiber.App, mqtt.Client, *redis.Client) {
	app := fiber.New()
	rdb, ctx := initRdb()
	broker := initBroker()
	initPostSubscription(broker, rdb)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api := app.Group("/api/v1")

	api.Post("/posts", func(c *fiber.Ctx) error {
		post := &Post{}
		err := post.Publish(broker, c)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.SendStatus(fiber.StatusOK)
	})

	api.Get("/posts", func(c *fiber.Ctx) error {
		keys, err := rdb.Keys(ctx, "*").Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.SendString(strings.Join(keys, ","))
	})

	api.Get("/search", func(c *fiber.Ctx) error {
		all, msgs, searched, err := searchPosts(rdb, c.Query("search", ""))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		jsonData, err := json.Marshal(msgs)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.SendString(fmt.Sprintf(`{ "total_keys": %d, "searched": %d, "messages": %s }`, all, searched, jsonData))
	})

	return app, broker, rdb
}
