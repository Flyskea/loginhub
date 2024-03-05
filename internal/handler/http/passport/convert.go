package passport

import (
	"fmt"
	"time"

	"github.com/mileusna/useragent"

	apiv1 "loginhub/api/v1/passport"
	entitypassport "loginhub/internal/domain/passport/entity"
	"loginhub/internal/domain/user/entity"
)

const (
	UserAgentKey = "User-Agent"
)

func UserEntityToVO(user *entity.User) *apiv1.User {
	return &apiv1.User{
		UserID: user.UserID,
		Name:   user.Name,
		Avatar: user.Avatar,
	}
}

func DevicesEntityToVO(device *entitypassport.LoginDevice) *apiv1.LoginDevice {
	devices := make([]*apiv1.Device, len(device.Devices))
	for i, d := range device.Devices {
		devices[i] = &apiv1.Device{
			ID:        d.ID,
			OS:        d.OS,
			Browser:   d.Browser,
			IP:        d.IP,
			Location:  d.Location,
			CreatedAt: d.CreatedAt,
		}
	}
	return &apiv1.LoginDevice{Devices: devices}
}

func NewDeviceEntity(userAgent string, ip string) *entitypassport.Device {
	agent := useragent.Parse(userAgent)

	return &entitypassport.Device{
		OS:        fmt.Sprintf("%s/%s", agent.OS, agent.OSVersion),
		Browser:   fmt.Sprintf("%s/%s", agent.Name, agent.Version),
		IP:        ip,
		CreatedAt: time.Now(),
	}
}
