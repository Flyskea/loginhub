package po

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type LoginDevice struct {
	ID        int64                 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	UserID    int64                 `gorm:"column:user_id;type:bigint;NOT NULL"`
	DeviceID  []byte                `gorm:"column:device_id;type:binary(16);NOT NULL"`
	Browser   string                `gorm:"column:browser;type:varchar(20);NOT NULL"`
	OS        string                `gorm:"column:os;type:varchar(20);NOT NULL"`
	IP        string                `gorm:"column:ip;type:varchar(40);NOT NULL"`
	CreatedAt time.Time             `gorm:"column:created_at;type:datetime(3);NOT NULL"`
	UpdatedAt time.Time             `gorm:"column:updated_at;type:datetime(3);NOT NULL"`
	DeletedAt soft_delete.DeletedAt `gorm:"column:deleted_at;type:bigint;default:0;NOT NULL"`
}

func (LoginDevice) TableName() string {
	return "login_device"
}
