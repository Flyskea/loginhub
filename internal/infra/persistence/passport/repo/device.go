package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Flyskea/gotools/errors"
	"github.com/huandu/go-sqlbuilder"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"loginhub/internal/base/iface"
	"loginhub/internal/base/reason"
	"loginhub/internal/conf"
	"loginhub/internal/domain/passport/entity"
	"loginhub/internal/domain/passport/repository"
	"loginhub/internal/infra/data"
	"loginhub/internal/infra/persistence/passport/convert"
	"loginhub/internal/infra/persistence/passport/po"
	commonconvert "loginhub/pkg/convert"
)

const (
	loginDevicePrefix = "device"
)

func loginDeviceKey(userID int64) string {
	return fmt.Sprintf("%s:%d", loginDevicePrefix, userID)
}

var _ repository.LoginDeviceRepository = (*LoginDeviceRepo)(nil)

type LoginDeviceRepo struct {
	txm  *data.TXManager
	rdb  redis.UniversalClient
	iprs iface.IP2RegionSearcher
}

func NewLoginDeviceRepo(
	txm *data.TXManager,
	rdb redis.UniversalClient,
	iprs iface.IP2RegionSearcher,
) *LoginDeviceRepo {
	return &LoginDeviceRepo{
		txm:  txm,
		rdb:  rdb,
		iprs: iprs,
	}
}

func (r *LoginDeviceRepo) fillLocation(ctx context.Context, device *entity.Device) error {
	region, err := r.iprs.RegionBasicByIPV4(ctx, device.IP)
	if err != nil {
		return err
	}
	device.Location = region.String()
	return nil
}

func (r *LoginDeviceRepo) fillLocations(ctx context.Context, device *entity.LoginDevice) error {
	for _, device := range device.Devices {
		if err := r.fillLocation(ctx, device); err != nil {
			return err
		}
	}
	return nil
}

func (r *LoginDeviceRepo) queryFn(ctx context.Context, userID int64) ([]*po.LoginDevice, error) {
	loginDevicePO := make([]*po.LoginDevice, 0)
	sb := sqlbuilder.NewSelectBuilder()
	sql, args := sb.From("login_device").
		Select("id", "user_id", "device_id", "os", "browser", "ip", "created_at").
		Where(sb.EQ("user_id", userID), sb.EQ("deleted_at", 0)).Build()
	err := r.txm.DB(ctx).Raw(sql, args...).Scan(&loginDevicePO).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFound(reason.DeviceNotExist)
		}
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if len(loginDevicePO) == 0 {
		return nil, errors.NotFound(reason.DeviceNotExist)
	}
	return loginDevicePO, nil
}

func (r *LoginDeviceRepo) cacheFn(ctx context.Context, userID int64, devices []*po.LoginDevice) error {
	if len(devices) == 0 {
		return nil
	}
	key := loginDeviceKey(userID)
	keyvalues := make([]interface{}, 0, len(devices)*2)
	for _, v := range devices {
		bytes, err := json.Marshal(v)
		if err != nil {
			return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		keyvalues = append(keyvalues, commonconvert.BytesToUUID(v.DeviceID), bytes)
	}

	err := r.rdb.HSet(ctx, key, keyvalues...).Err()
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}

	err = r.rdb.Expire(ctx, key, conf.GConf.GetPassport().GetDeviceCacheTtl().AsDuration()).Err()
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

func (r *LoginDeviceRepo) GetByUserID(ctx context.Context, userID int64) (*entity.LoginDevice, error) {
	key := loginDeviceKey(userID)
	loginDevicePO := make([]*po.LoginDevice, 0)

	loginDevices, err := r.rdb.HGetAll(ctx, key).Result()
	switch {
	case err == redis.Nil || len(loginDevices) == 0:
		loginDevicePO, err = r.queryFn(ctx, userID)
		if err != nil {
			return nil, err
		}
		err = r.cacheFn(ctx, userID, loginDevicePO)
		if err != nil {
			return nil, err
		}
	case err != nil:
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	default:
		for _, v := range loginDevices {
			devicePO := po.LoginDevice{}
			err := json.Unmarshal(commonconvert.StringToBytes(v), &devicePO)
			if err != nil {
				return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
			}
			loginDevicePO = append(loginDevicePO, &devicePO)
		}
	}

	loginDeviceEntity := convert.DevicePOsToEntity(loginDevicePO)
	err = r.fillLocations(ctx, loginDeviceEntity)
	if err != nil {
		return nil, err
	}
	return loginDeviceEntity, nil
}

func (r *LoginDeviceRepo) GetDeviceByID(ctx context.Context, userID int64, deviceID string) (*entity.Device, error) {
	key := loginDeviceKey(userID)
	var devicePO *po.LoginDevice
	value, err := r.rdb.HGet(ctx, key, deviceID).Result()

	switch {
	case err == redis.Nil || value == "":
		var loginDevicePO []*po.LoginDevice
		loginDevicePO, err = r.queryFn(ctx, userID)
		if err != nil {
			return nil, err
		}
		err = r.cacheFn(ctx, userID, loginDevicePO)
		if err != nil {
			return nil, err
		}

		for _, v := range loginDevicePO {
			if commonconvert.BytesToUUID(v.DeviceID) == deviceID {
				devicePO = v
				break
			}
		}
	case err != nil:
		return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	default:
		devicePO = &po.LoginDevice{}
		err := json.Unmarshal(commonconvert.StringToBytes(value), devicePO)
		if err != nil {
			return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
	}

	if devicePO == nil {
		return nil, errors.NotFound(reason.DeviceNotExist)
	}

	deviceEntity := convert.DevicePOToEntity(devicePO)
	err = r.fillLocation(ctx, deviceEntity)
	if err != nil {
		return nil, err
	}
	return deviceEntity, nil
}

func (r *LoginDeviceRepo) Create(ctx context.Context, userID int64, device *entity.Device) error {
	loginDevicePO := convert.DevicePOFromEntity(device)
	loginDevicePO.UserID = userID
	return r.txm.ExecTx(ctx, func(ctx context.Context) error {
		err := r.txm.DB(ctx).Create(&loginDevicePO).Error
		if err != nil {
			return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}

		err = r.cacheFn(ctx, userID, []*po.LoginDevice{loginDevicePO})
		if err != nil {
			return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		return nil
	})
}

func (r *LoginDeviceRepo) DeleteOneByDeviceID(ctx context.Context, userID int64, deviceID string) error {
	key := loginDeviceKey(userID)
	err := r.rdb.HDel(ctx, key, deviceID).Err()
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	sb := sqlbuilder.NewUpdateBuilder()
	sql, args := sb.Where(
		sb.Equal("user_id", userID),
		sb.EQ("device_id", commonconvert.BytesToString(commonconvert.UUIDToBytes(deviceID))),
	).Set(sb.Assign("deleted_at", time.Now().Unix())).
		Update("login_device").Build()
	fmt.Println(sql, args)
	err = r.txm.DB(ctx).Exec(sql, args...).Error
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}

func (r *LoginDeviceRepo) DeleteAllByUserID(ctx context.Context, userID int64) error {
	key := loginDeviceKey(userID)
	err := r.rdb.Del(ctx, key).Err()
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	sb := sqlbuilder.NewUpdateBuilder()
	sql, args := sb.Where(
		sb.Equal("user_id", userID),
	).Set(sb.Assign("deleted_at", time.Now().Unix())).
		Update("login_device").Build()
	err = r.txm.DB(ctx).Exec(sql, args...).Error
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return nil
}
