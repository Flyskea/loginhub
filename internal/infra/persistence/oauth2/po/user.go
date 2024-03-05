package po

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type UserThirdAuth struct {
	ID           int64                 `gorm:"column:id;primary_key"`
	UserID       int64                 `gorm:"column:user_id;NOT NULL"`
	Type         string                `gorm:"column:type;NOT NULL"`          // OAUTH2 provider type, e.g. github, wechat, etc.
	AuthID       string                `gorm:"column:auth_id;NOT NULL"`       // 第三方 uid 、openid 等
	UnionID      string                `gorm:"column:union_id;NOT NULL"`      // QQ / 微信同一主体下 Unionid 相同
	Credential   string                `gorm:"column:credential;NOT NULL"`    // access_token
	RefreshToken string                `gorm:"column:refresh_token;NOT NULL"` // refresh_token
	CreatedAt    time.Time             `gorm:"column:created_at;type:datetime;NOT NULL"`
	UpdatedAt    time.Time             `gorm:"column:updated_at;type:datetime;NOT NULL"`
	DeletedAt    soft_delete.DeletedAt `gorm:"column:deleted_at;type:bigint;NOT NULL"`
}

const (
	UserThirdAuthTableName = "user_third_auth"
)

func (UserThirdAuth) TableName() string {
	return UserThirdAuthTableName
}
