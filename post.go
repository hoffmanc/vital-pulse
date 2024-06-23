package main

import (
	"encoding/json"
	"log"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-redis/redis/v8"
	fiber "github.com/gofiber/fiber/v2"
	"golang.org/x/net/context"
)

type Post struct {
	Id         string    `json:"id"`
	Body       string    `json:"body"`
	From       string    `json:"from"`
	Subject    string    `json:"subject"`
	ReceivedAt GmailTime `json:"receivedAt"`
}

func (p *Post) Publish(broker mqtt.Client, c *fiber.Ctx) error {
	t := broker.Publish("inbox", 0, false, c.Body())
	// go func() {
	<-t.Done()
	if t.Error() != nil {
		return t.Error()
	}
	// }()

	return nil
}

func initPostSubscription(broker mqtt.Client, rdb *redis.Client) {
	broker.Subscribe("inbox", 0, savePost(rdb))
}

func savePost(rdb *redis.Client) func(client mqtt.Client, msg mqtt.Message) {
	return func(client mqtt.Client, msg mqtt.Message) {
		var p Post
		err := json.Unmarshal(msg.Payload(), &p)
		if err != nil {
			log.Println(err)
			return
		}

		if err := rdb.Set(context.TODO(), p.Id, msg.Payload(), 0).Err(); err != nil {
			log.Println(err)
			return
		}
	}
}

func searchPosts(rdb *redis.Client, search string) (int, []Post, int, error) {
	posts := []Post{}

	ctx := context.Background()
	keys, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		log.Println("rdb.Keys", err)
		return len(keys), posts, 0, err
	}

	var i int64 = 0
	for {
		var batch int64 = 100
		if len(keys[i:]) < int(batch) {
			batch = int64(len(keys[i:]))
		}
		log.Println("batch", batch)

		msgJSONs, err := rdb.MGet(ctx, keys[i:i+batch]...).Result()
		if err != nil {
			log.Println("rdb.MGet", err)
			return len(keys), posts, 0, err
		}

		for _, msgJSON := range msgJSONs {
			var post Post
			s := msgJSON.(string)
			log.Println("s", s)
			err := json.Unmarshal([]byte(s), &posts)
			if err != nil {
				log.Println("json.Unmarshal", err, s)
			}

			if strings.Contains(post.Body, search) {
				posts = append(posts, post)
			}
		}

		if batch < 100 {
			i += batch
			break
		}
		i += 100
	}
	return len(keys), posts, int(i), nil
}
