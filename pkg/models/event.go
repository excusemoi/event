package models

import "time"

type Event struct {
	Id     uint64    `json:"id"`
	UserId uint64    `json:"user_id"`
	Name   string    `json:"name"`
	Date   time.Time `json:"date"`
}
