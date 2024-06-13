package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	fiber "github.com/gofiber/fiber/v2"
	"golang.org/x/net/context"
)

type Message struct {
	Id         string    `json:"id"`
	Body       string    `json:"body"`
	From       string    `json:"from"`
	Subject    string    `json:"subject"`
	ReceivedAt GmailTime `json:"receivedAt"`
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
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		date := &msg.ReceivedAt
		key := fmt.Sprintf("%s:%s", date.Format(time.RFC3339), msg.Id)
		err = rdb.SetEX(ctx, key, msg.Body, 24*time.Hour).Err()
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.SendStatus(fiber.StatusOK)
	})

	api.Get("/posts", func(c *fiber.Ctx) error {
		ctx := context.Background()
		keys, err := rdb.Keys(ctx, "2024*").Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.SendString(strings.Join(keys, ","))
	})

	api.Get("/search", func(c *fiber.Ctx) error {
		msgs, err := searchPosts(rdb, c.Query("search", ""))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		jsonData, err := json.Marshal(msgs)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.SendString(string(jsonData))
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

func searchPosts(rdb *redis.Client, s string) ([]Message, error) {
	msgs := []Message{}

	ptn, err := regexp.Compile(s)
	if err != nil {
		log.Println("regexp.Compile", err)
		return msgs, err
	}

	ctx := context.Background()
	keys, err := rdb.Keys(ctx, "2024*").Result()
	if err != nil {
		log.Println("rdb.Keys", err)
		return msgs, err
	}

	var i int64 = 0
	for {
		var batch int64 = 100
		if len(keys[i:]) < int(batch) {
			batch = int64(len(keys[i:]))
		}
		msgJSONs, err := rdb.MGet(ctx, keys[i:i+batch]...).Result()
		if err != nil {
			log.Println("rdb.MGet", err)
			return msgs, err
		}

		for _, msgJSON := range msgJSONs {
			var msg Message
			s := msgJSON.(string)
			err := json.Unmarshal([]byte(s), &msg)
			if err != nil {
				log.Println("json.Unmarshal", err, s)
			}

			for _, m := range msgs {
				if ptn.MatchString(m.Body) {
					msgs = append(msgs, m)
				}
			}
		}

		if batch < 100 {
			break
		}
		i += 100
	}
	return msgs, nil
}
