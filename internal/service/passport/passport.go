package passport

import (
	"context"
	"time"

	"github.com/Flyskea/gotools/errors"

	"loginhub/internal/base/iface"
	"loginhub/internal/base/reason"
	"loginhub/internal/domain/passport/entity"
	"loginhub/internal/domain/passport/repository"
	"loginhub/internal/domain/passport/service"
	userentity "loginhub/internal/domain/user/entity"
	userrepo "loginhub/internal/domain/user/repository"
	"loginhub/internal/service/captcha"
)

type PassportService struct {
	tx  iface.Transaction
	pds *service.PassportService
	pdr repository.LoginDeviceRepository
	atr repository.AccessTokenRepository
	rtr repository.RefreshTokenRepository
	ur  userrepo.UserRepository
	uig iface.UniqueIDGenerator
	cs  *captcha.CaptchaService
}

func NewPassportService(
	tx iface.Transaction,
	pds *service.PassportService,
	pdr repository.LoginDeviceRepository,
	atr repository.AccessTokenRepository,
	rtr repository.RefreshTokenRepository,
	ur userrepo.UserRepository,
	uig iface.UniqueIDGenerator,
	cs *captcha.CaptchaService,
) *PassportService {
	return &PassportService{
		tx:  tx,
		pds: pds,
		pdr: pdr,
		atr: atr,
		rtr: rtr,
		ur:  ur,
		uig: uig,
		cs:  cs,
	}
}

func (s *PassportService) EmailCaptcha(ctx context.Context, e *EmailSend) error {
	_, err := s.cs.EmailCaptcha(ctx, e.Email, e.CaptchaType.ToCaptchaService())
	if err != nil {
		return err
	}
	return nil
}

func (s *PassportService) registerByEmail(ctx context.Context, req *RegisterInfo) (*userentity.User, error) {
	err := s.cs.IsCodeCorrectByEmail(ctx, req.Email, captcha.RegisterCaptchatType, req.Code)
	if err != nil {
		return nil, err
	}

	_, err = s.ur.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.BadRequest(reason.EmailDuplicate)
	} else if !errors.IsNotFound(err) {
		return nil, err
	}
	uniqueID, err := s.uig.NextIDInt64(ctx)
	if err != nil {
		return nil, err
	}
	user := userentity.NewUser(
		uniqueID,
		req.UserName,
		req.Password,
		req.Email,
		"",
		"")
	user.ActiveInfo = &userentity.ActiveInfo{
		LastLoginAt: time.Now(),
		IP:          req.Device.IP,
	}

	return user, nil
}

func (s *PassportService) Register(ctx context.Context, req *RegisterInfo) (*LoginResult, error) {
	if req.Password != req.PasswordConfirm {
		return nil, errors.BadRequest(reason.PasswordNotMatch)
	}

	var (
		user *userentity.User
		err  error
	)
	switch req.RegisterType {
	case EmailRegisterType:
		user, err = s.registerByEmail(ctx, req)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.BadRequest(reason.RegisterTypeNotSupport)
	}

	authn, err := s.pds.Register(ctx, &service.RegisterInfo{
		RegisterType: service.RegisterType(req.RegisterType),
		User:         user,
		Device:       req.Device,
	})
	if err != nil {
		return nil, err
	}
	user = authn.User
	var signedAccessToken, signedRefreshToken string
	err = s.tx.ExecTx(ctx, func(ctx context.Context) error {
		err = s.ur.Create(ctx, user)
		if err != nil {
			return err
		}
		err = s.pdr.Create(ctx, user.UserID, req.Device)
		if err != nil {
			return err
		}
		err = s.atr.StoreToken(ctx, authn.AccessToken)
		if err != nil {
			return err
		}
		err = s.rtr.StoreToken(ctx, authn.RefreshToken)
		if err != nil {
			return err
		}
		signedAccessToken, err = authn.AccessToken.Sign()
		if err != nil {
			return err
		}
		signedRefreshToken, err = authn.AccessToken.Sign()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &LoginResult{
		AccessToken:    signedAccessToken,
		AccessTokenTTL: int(authn.AccessToken.GetTTL()),
		RefreshToken:   signedRefreshToken,
		User:           user,
	}, nil
}

func (s *PassportService) Login(ctx context.Context, info *service.LoginInfo) (*LoginResult, error) {
	resp, err := s.pds.Authentication(ctx, info)
	if err != nil {
		return nil, err
	}

	var (
		refreshToken string
		accessToken  string
	)
	err = s.tx.ExecTx(ctx, func(ctx context.Context) error {
		err = s.pdr.Create(ctx, resp.User.UserID, info.Device)
		if err != nil {
			return err
		}
		err = s.atr.StoreToken(ctx, resp.AccessToken)
		if err != nil {
			return err
		}
		err = s.rtr.StoreToken(ctx, resp.RefreshToken)
		if err != nil {
			return err
		}
		accessToken, err = resp.AccessToken.Sign()
		if err != nil {
			return err
		}
		refreshToken, err = resp.RefreshToken.Sign()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &LoginResult{
		AccessToken:    accessToken,
		AccessTokenTTL: int(resp.AccessToken.ExpiresAt.Sub(resp.AccessToken.IssuedAt.Time).Seconds()),
		RefreshToken:   refreshToken,
		User:           resp.User,
	}, nil
}

func (s *PassportService) Logout(ctx context.Context, userInfo *entity.UserInfo) error {
	return s.deleteDeviceCore(ctx, userInfo.UserID, userInfo.DeviceID)
}

func (s *PassportService) Refresh(ctx context.Context, refreshToken string) (*RefreshResult, error) {
	token, err := entity.ParseAndVerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	isExist, err := s.rtr.IsTokenExist(ctx, token)
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, errors.Unauthorized(reason.UnauthorizedError)
	}
	resp, err := s.pds.Refresh(ctx, token)
	if err != nil {
		return nil, err
	}
	var accessToken string
	err = s.tx.ExecTx(ctx, func(ctx context.Context) error {
		err = s.rtr.RotateToken(ctx, token, resp.RefreshToken)
		if err != nil {
			return err
		}
		err = s.atr.StoreToken(ctx, resp.AccessToken)
		if err != nil {
			return err
		}
		accessToken, err = resp.AccessToken.Sign()
		if err != nil {
			return err
		}
		refreshToken, err = resp.RefreshToken.Sign()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &RefreshResult{
		AccessToken:    accessToken,
		AccessTokenTTL: int(resp.AccessToken.ExpiresAt.Sub(resp.AccessToken.IssuedAt.Time).Seconds()),
		RefreshToken:   refreshToken,
		User:           resp.User,
	}, nil
}

func (s *PassportService) LoginDevices(ctx context.Context, userID int64) (*entity.LoginDevice, error) {
	return s.pdr.GetByUserID(ctx, userID)
}

func (s *PassportService) deleteDeviceCore(ctx context.Context, userID int64, deviceID string) error {
	return s.tx.ExecTx(ctx, func(ctx context.Context) error {
		err := s.pdr.DeleteOneByDeviceID(ctx, userID, deviceID)
		if err != nil {
			return err
		}
		err = s.atr.RevokeTokenByDeviceID(ctx, userID, deviceID)
		if err != nil {
			return err
		}
		err = s.rtr.RevokeTokenByDeviceID(ctx, userID, deviceID)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *PassportService) DeleteDevice(ctx context.Context, req *DeleteDevice) error {
	if req.UserInfo.DeviceID == req.DeviceID {
		return errors.BadRequest(reason.CannotDeleteCurrentDevice)
	}
	_, err := s.pdr.GetDeviceByID(ctx, req.UserInfo.UserID, req.DeviceID)
	if err != nil {
		return err
	}

	return s.deleteDeviceCore(ctx, req.UserInfo.UserID, req.DeviceID)
}

func (s *PassportService) VerifyAccessToken(ctx context.Context, accessToken string) (*entity.AccessToken, error) {
	token, err := entity.ParseAndVerifyAccessToken(accessToken)
	if err != nil {
		return nil, err
	}
	isExist, err := s.atr.IsTokenExist(ctx, token)
	if err != nil {
		return nil, err
	}
	if !isExist {
		return nil, errors.Unauthorized(reason.UnauthorizedError)
	}
	return token, nil
}
