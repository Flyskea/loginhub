package entity

import "time"

type ActiveInfo struct {
	LastLoginAt time.Time
	IP     string // user current ip
}
