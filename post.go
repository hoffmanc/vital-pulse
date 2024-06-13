package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

func searchPosts(rdb *redis.Client, s string) ([]Message, int, error) {
	msgs := []Message{}

	ptn, err := regexp.Compile(fmt.Sprintf(".*%s.*", s))
	if err != nil {
		log.Println("regexp.Compile", err)
		return msgs, 0, err
	}

	ctx := context.Background()
	keys, err := rdb.Keys(ctx, "2024*").Result()
	if err != nil {
		log.Println("rdb.Keys", err)
		return msgs, 0, err
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
			return msgs, 0, err
		}

		for _, msgJSON := range msgJSONs {
			var msg Message
			s := msgJSON.(string)
			err := json.Unmarshal([]byte(s), &msg)
			if err != nil {
				log.Println("json.Unmarshal", err, s)
			}

			if ptn.MatchString(msg.Body) {
				msgs = append(msgs, msg)
			}
		}

		if batch < 100 {
			break
		}
		i += 100
	}
	return msgs, int(i), nil
}
