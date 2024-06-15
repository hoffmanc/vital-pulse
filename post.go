package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

func searchPosts(rdb *redis.Client, search string) (int, []Message, int, error) {
	msgs := []Message{}

	ctx := context.Background()
	keys, err := rdb.Keys(ctx, "*").Result()
	if err != nil {
		log.Println("rdb.Keys", err)
		return len(keys), msgs, 0, err
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
			return len(keys), msgs, 0, err
		}

		for _, msgJSON := range msgJSONs {
			var msg Message
			s := msgJSON.(string)
			log.Println("s", s)
			err := json.Unmarshal([]byte(s), &msg)
			if err != nil {
				log.Println("json.Unmarshal", err, s)
			}

			if strings.Contains(msg.Body, search) {
				msgs = append(msgs, msg)
			}
		}

		if batch < 100 {
			i += batch
			break
		}
		i += 100
	}
	return len(keys), msgs, int(i), nil
}
