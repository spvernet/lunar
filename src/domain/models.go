package domain

import "time"

type RocketStatus string

const (
	StatusActive   RocketStatus = "ACTIVE"
	StatusExploded RocketStatus = "EXPLODED"
)

type Rocket struct {
	Channel    string
	Type       string
	Mission    string
	Speed      int64
	Status     RocketStatus
	LastMsgNum int
	UpdatedAt  time.Time
}
