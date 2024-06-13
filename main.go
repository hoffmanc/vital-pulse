package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	fiber "github.com/gofiber/fiber/v2"
	"golang.org/x/net/context"
)

type Message struct {
	Id         string    `json:"id"`
	Body       string    `json:"body"`
	Subject    string    `json:"subject"`
	ReceivedAt time.Time `json:"receivedAt"`
}

func main() {
	app := fiber.New()
	rdb, ctx := initRdb()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api := app.Group("/api/v1")

	api.Post("/posts", func(c *fiber.Ctx) error {
		msg := Message{}
		err := json.Unmarshal(c.Body(), &msg)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		date := &msg.ReceivedAt
		key := fmt.Sprintf("%s:%s", date.Format(time.UnixDate), msg.Id)
		err = rdb.SetEX(ctx, key, msg.Body, 24*time.Hour).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.SendStatus(fiber.StatusOK)
	})

	api.Get("/posts", func(c *fiber.Ctx) error {
		ctx := context.Background()
		keys := rdb.Keys(ctx, "2024*")

		return c.SendString(strings.Join(keys))
	})

	port := os.Getenv("PORT")
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}

func initRdb() (*redis.Client, context.Context) {
	// Parse the URL
	u, err := url.Parse(os.Getenv("REDISCLOUD_URL"))
	if err != nil {
		log.Fatal(err)
	}

	pwd, _ := u.User.Password()

	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     u.Host, // also has port for some reason
		Username: u.User.Username(),
		Password: pwd,
	})

	ctx := context.Background()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	return rdb, ctx
}
