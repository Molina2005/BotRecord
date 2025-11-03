package models

import "time"

type Recordatorio struct {
	Id         int
	Title      string
	DateRecord time.Time
}
