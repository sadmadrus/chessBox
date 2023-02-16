package model

import "time"

type Session struct {
	Uid          uint64
	StartDate    time.Time
	EndDate      time.Time
	BlackPlayer  uint64
	WhitePlayer  uint64
	Description  string
	SessionState uint8
	Position     string
}
