package main

import (
	"fmt"
	"strings"
	"time"
)

const gmailLayout = "Date: Thu, 02 Jan 2006 15:04:05 -0700"

type GmailTime struct {
	time.Time
}

func (t *GmailTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		t.Time = time.Time{}
		return
	}
	t.Time, err = time.Parse(gmailLayout, s)
	return
}

func (t *GmailTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", t.Time.Format(gmailLayout))), nil
}

func (t *GmailTime) IsSet() bool {
	return !t.IsZero()
}
