package models

import "time"

type Recordatorio struct {
	IdRecordatorios int64
	Id              int64
	Title           string
	DateRecord      time.Time
	Estado          string
}
