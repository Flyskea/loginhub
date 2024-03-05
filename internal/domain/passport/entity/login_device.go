package entity

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID        string
	OS        string
	Browser   string
	IP        string
	Location  string
	CreatedAt time.Time
}

func NewDevice() *Device {
	id := uuid.New().String()
	return &Device{
		ID:        id,
		CreatedAt: time.Now(),
	}
}

type LoginDevice struct {
	UserID  int64
	Devices []*Device
}
