package models

import "time"

type ClickEvents struct {
	ShortKey      string      `json:"short_key"`
	TimeStamp     time.Time   `json:"timestamp"`
	IP            string      `json:"ip"`
	Country       string      `json:"country"`
	Browser       string      `json:"browser"`
	Device        string      `json:"device"`
	Referrer      string      `json:"referrer"`
	UserAgent     string      `json:"user_agent"`
}
