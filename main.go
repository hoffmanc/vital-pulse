package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/go-redis/redis/v8"
	fiber "github.com/gofiber/fiber/v2"
	"golang.org/x/net/context"
)

func main() {
	app := fiber.New()
	rdb, ctx := initRdb()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api := app.Group("/api/v1")

	api.Post("/posts", func(c *fiber.Ctx) error {
		err := rdb.Set(ctx, "email", c.Body(), 0).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.SendStatus(fiber.StatusOK)
	})

	api.Get("/posts", func(c *fiber.Ctx) error {
		post, err := rdb.Get(ctx, "email").Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.SendString(post)
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

	log.Println("u.Host", u.Host)
	log.Println("u.User.String()", u.User.String())
	log.Println("u.User.Password()", pwd)

	// Create Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     u.Host, // also has port for some reason
		Username: u.User.String(),
		Password: pwd,
	})

	ctx := context.Background()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	return rdb, ctx
}
