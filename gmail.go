package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type GmailTime struct {
	time.Time
}

func (t *GmailTime) UnmarshalJSON(b []byte) error {
	gmailLayouts := []string{
		"Date: Mon, 02 Jan 2006 15:04:05 -0700",
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"Date: Mon, 02 Jan 2006 15:04:05 MST",
		"Mon, 02 Jan 2006 15:04:05 MST",
	}

	s := strings.Trim(string(b), "\"")
	if s == "null" {
		t.Time = time.Time{}
		return nil
	}
	var err error

	for _, layout := range gmailLayouts {
		t.Time, err = time.Parse(layout, s)
		if err == nil {
			log.Printf("Successfully parsed %s to %s", s, t.Time)
			return nil
		}
		log.Println(err)
	}
	return err
}

func (t *GmailTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", t.Time.Format("Mon, 02 Jan 2006 15:04:05 MST"))), nil
}

func (t *GmailTime) IsSet() bool {
	return !t.IsZero()
}
