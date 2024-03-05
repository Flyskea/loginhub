package po

import (
	"database/sql"
	"time"

	"gorm.io/plugin/soft_delete"
)

type User struct {
	ID          int64                 `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT"`
	UserID      int64                 `gorm:"column:user_id;type:bigint(20);NOT NULL"`
	Name        string                `gorm:"column:name;type:varchar(20);NOT NULL"`
	Password    string                `gorm:"column:password;type:varchar(50);NOT NULL"`
	Avatar      string                `gorm:"column:avatar;type:varchar(255);NOT NULL"`
	Mobile      sql.NullString        `gorm:"column:mobile;type:varchar(20)"`
	Email       sql.NullString        `gorm:"column:email;type:varchar(255)"`
	Account     sql.NullString        `gorm:"column:account;type:varchar(255)"`
	LastLoginAt int64                 `gorm:"column:last_login_at;type:bigint(20);NOT NULL"` // 最后登陆时间
	IP          string                `gorm:"column:ip;type:varchar(45);NOT NULL"`
	CreatedAt   time.Time             `gorm:"column:created_at;type:datetime;NOT NULL"`
	UpdatedAt   time.Time             `gorm:"column:updated_at;type:datetime;NOT NULL"`
	DeletedAt   soft_delete.DeletedAt `gorm:"column:deleted_at;type:bigint;NOT NULL"`
}

const (
	UserTableName = "user"
)

func (User) TableName() string {
	return UserTableName
}
