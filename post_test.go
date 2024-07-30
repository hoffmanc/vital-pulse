package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

// TestCreate tests the CreateMessage function.
func TestCreate(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	body := &Post{
		Id:         "test:a1",
		Body:       "here is the email body",
		Subject:    "an email subject",
		From:       "foo@bar.com",
		ReceivedAt: GmailTime{Time: time.Now()},
	}

	jsonBody, err := json.Marshal(body)
	assert.NoError(t, err)

	req := httptest.NewRequest(
		"POST",
		"/api/v1/posts",
		bytes.NewReader(jsonBody),
	)

	app, _, rdb := InitApp()
	_, err = rdb.Del(context.TODO(), body.Id).Result()
	assert.NoError(t, err)

	resp, _ := app.Test(req, 1)

	assert.Equalf(t, 200, resp.StatusCode, "nope")

	assert.Eventually(t, func() bool {
		result, err := rdb.Keys(context.TODO(), body.Id).Result()
		if err != nil {
			assert.NoError(t, err)
			return false
		}
		assert.NotEmptyf(t, result, "not empty")
		return true
	}, time.Second*10, time.Millisecond*100)

}
