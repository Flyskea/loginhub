package po

import (
	"time"
)

type OAuth2Provider struct {
	ID           int64     `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`
	Type         string    `gorm:"column:type;type:varchar(20);NOT NULL"` // OAUTH2 provider type, e.g. github, wechat, etc.
	ClientID     string    `gorm:"column:client_id;type:varchar(255);NOT NULL"`
	ClientSecret string    `gorm:"column:client_secret;type:varchar(255);NOT NULL"`
	RedirectUrl  string    `gorm:"column:redirect_url;type:varchar(255);NOT NULL"`
	CreatedAt    time.Time `gorm:"column:created_at;type:datetime;NOT NULL"`
	UpdatedAt    time.Time `gorm:"column:updated_at;type:datetime;NOT NULL"`
	DeletedAt    int64     `gorm:"column:deleted_at;type:bigint;NOT NULL"`
}

func (OAuth2Provider) TableName() string {
	return OAuth2ProviderTableName
}

const (
	OAuth2ProviderTableName = "oauth2_provider"
)

type OAuth2ProviderList []*OAuth2Provider
