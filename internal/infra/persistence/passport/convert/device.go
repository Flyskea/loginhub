package convert

import (
	"loginhub/internal/domain/passport/entity"
	"loginhub/internal/infra/persistence/passport/po"
	"loginhub/pkg/convert"
)

func DevicePOFromEntity(device *entity.Device) *po.LoginDevice {
	return &po.LoginDevice{
		DeviceID:  convert.UUIDToBytes(device.ID),
		OS:        device.OS,
		Browser:   device.Browser,
		IP:        device.IP,
		CreatedAt: device.CreatedAt,
		UpdatedAt: device.CreatedAt,
	}
}

func DevicePOToEntity(device *po.LoginDevice) *entity.Device {
	return &entity.Device{
		ID:        convert.BytesToUUID(device.DeviceID),
		OS:        device.OS,
		Browser:   device.Browser,
		IP:        device.IP,
		CreatedAt: device.CreatedAt,
	}
}

func DevicePOsToEntity(devices []*po.LoginDevice) *entity.LoginDevice {
	var entities []*entity.Device
	for _, device := range devices {
		entities = append(entities, DevicePOToEntity(device))
	}
	return &entity.LoginDevice{
		UserID:  devices[0].UserID,
		Devices: entities,
	}
}
