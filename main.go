package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	fiber "github.com/gofiber/fiber/v2"
	"golang.org/x/net/context"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api := app.Group("/api/v1")

	api.Post("/posts", func(c *fiber.Ctx) error {
		// err := rdb.Set(ctx, "email", c.Body(), 0).Err()
		// if err != nil {
		// 	return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		// }

		return c.SendStatus(fiber.StatusOK)
	})

	api.Get("/posts", func(c *fiber.Ctx) error {
		// post, err := rdb.Get(ctx, "email").Result()
		// if err != nil {
		// 	return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		// }

		return c.SendString("foop")
	})

	port := os.Getenv("PORT")
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}

func initRdb() *redis.Client {
	// Read Redis connection details from environment variables
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	return rdb
}
