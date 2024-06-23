package main

import (
	"context"
	"log"
	"net/url"
	"os"

	"github.com/go-redis/redis/v8"
)

func initRdb() (*redis.Client, context.Context) {
	// Parse the URL
	u, err := url.Parse(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal(err)
	}

	// Create Redis client
	opts := &redis.Options{
		Addr: u.Host, // also has port for some reason
	}

	if u.User.Username() != "" {
		opts.Username = u.User.Username()
	}

	pwd, _ := u.User.Password()
	if pwd != "" {
		opts.Password = pwd
	}

	rdb := redis.NewClient(opts)

	ctx := context.Background()
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	return rdb, ctx
}
