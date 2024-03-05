package passport

import (
	"context"
	"time"

	"github.com/Flyskea/gotools/errors"

	"loginhub/internal/base/iface"
	"loginhub/internal/base/reason"
	oauth2repo "loginhub/internal/domain/oauth2/repository"
	"loginhub/internal/domain/passport/entity"
	"loginhub/internal/domain/passport/repository"
	"loginhub/internal/domain/passport/service"
	userentity "loginhub/internal/domain/user/entity"
	userrepo "loginhub/internal/domain/user/repository"
	"loginhub/internal/service/captcha"
	"loginhub/internal/service/oauth2"
)

const (
	defaultUserName    = "admin"
	defaultPassword    = "admin"
	defaultEmail       = "admin@localhost"
	defaultMobile      = "16666"
	defaultAccountName = "admin"
)

type PassportService struct {
	tx            iface.Transaction
	pds           *service.PassportService
	pdr           repository.LoginDeviceRepository
	atr           repository.AccessTokenRepository
	rtr           repository.RefreshTokenRepository
	our           oauth2repo.OAuth2UserRepository
	ur            userrepo.UserRepository
	uig           iface.UniqueIDGenerator
	cs            *captcha.CaptchaService
	oas           *oauth2.OAuth2Service
	userSearchers []LocalUserSearcher
}

func NewPassportService(
	tx iface.Transaction,
	pds *service.PassportService,
	pdr repository.LoginDeviceRepository,
	atr repository.AccessTokenRepository,
	rtr repository.RefreshTokenRepository,
	our oauth2repo.OAuth2UserRepository,
	ur userrepo.UserRepository,
	uig iface.UniqueIDGenerator,
	oas *oauth2.OAuth2Service,
	cs *captcha.CaptchaService,
) (*PassportService, error) {
	userSearchers := make([]LocalUserSearcher, 0, 3)
	userSearchers = append(userSearchers, NewEmailAccountSearcher(ur))
	userSearchers = append(userSearchers, NewAccountSearcher(ur))
	userSearchers = append(userSearchers, NewMobileAccountSearcher(ur))
	s := &PassportService{
		tx:            tx,
		pds:           pds,
		pdr:           pdr,
		atr:           atr,
		rtr:           rtr,
		our:           our,
		ur:            ur,
		uig:           uig,
		cs:            cs,
		oas:           oas,
		userSearchers: userSearchers,
	}
	return s, s.Init(context.Background())
}

func (s *PassportService) EmailCaptcha(ctx context.Context, e *EmailSend) error {
	_, err := s.cs.EmailCaptcha(ctx, e.Email, e.CaptchaType.ToCaptchaService())
	if err != nil {
		return err
	}
	return nil
}

func (s *PassportService) registerUser(
	ctx context.Context,
	username string,
	password string,
	email string,
	mobile string,
	avatar string,
	device *entity.Device,
) (*userentity.User, error) {
	uniqueID, err := s.uig.NextIDInt64(ctx)
	if err != nil {
		return nil, err
	}
	user := userentity.NewUser(
		uniqueID,
		username,
		password,
		email,
		mobile,
		avatar)
	user.ActiveInfo = &userentity.ActiveInfo{
		LastLoginAt: time.Now(),
		IP:          device.IP,
	}
	return user, nil
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

	return s.registerUser(ctx, req.UserName, req.Password, req.Email,
		"", "", req.Device)
}

func (s *PassportService) Register(ctx context.Context, req *RegisterInfo) (*AuthenticationResult, error) {
	if req.Password != req.PasswordConfirm {
		return nil, errors.BadRequest(reason.PasswordMismatch)
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
	return &AuthenticationResult{
		AccessToken:    signedAccessToken,
		AccessTokenTTL: int(authn.AccessToken.GetTTL()),
		RefreshToken:   signedRefreshToken,
		User:           user,
	}, nil
}

func (s *PassportService) oauth2Login(ctx context.Context, l *LoginAction) (*service.AuthenticationInfo, error) {
	result, err := s.oas.OAuth2UserInfo(ctx, &oauth2.GetUserInfoAction{
		ProviderType:    l.Provider,
		Code:            l.Code,
		RequestState:    l.RequestState,
		OAuth2SessionID: l.OAuth2SessionID,
	})
	if err != nil {
		return nil, err
	}
	var resp *service.AuthenticationInfo
	user, err := s.our.GetByOAuth2UserID(ctx, result.UserInfo.ID)
	switch {
	case err == nil:
		resp, err = s.pds.Authenticate(ctx, &service.AuthenticationParams{
			Password: "",
			User:     user,
			Device:   l.Device,
		})
		if err != nil {
			return nil, err
		}
	case errors.IsNotFound(err):
		user, err = s.ur.GetByEmail(ctx, result.UserInfo.Email)
		if err != nil && !errors.IsNotFound(err) {
			return nil, err
		}
		if user != nil {
			err = s.our.Save(ctx, user.UserID, l.Provider, result.Token, result.UserInfo)
			if err != nil {
				return nil, err
			}
		}
		if user == nil {
			user, err = s.ur.GetByMobile(ctx, result.UserInfo.Phone)
			if err != nil && !errors.IsNotFound(err) {
				return nil, err
			}
			if user != nil {
				err = s.our.Save(ctx, user.UserID, l.Provider, result.Token, result.UserInfo)
				if err != nil {
					return nil, err
				}
			}
		}
		if user == nil {
			user, err = s.registerUser(ctx, result.UserInfo.Username, "", result.UserInfo.Email,
				result.UserInfo.Phone, result.UserInfo.AvatarUrl, l.Device)
			if err != nil {
				return nil, err
			}
			err = s.tx.ExecTx(ctx, func(ctx context.Context) error {
				err = s.our.Save(ctx, user.UserID, l.Provider, result.Token, result.UserInfo)
				if err != nil {
					return err
				}
				return s.ur.Create(ctx, user)
			})
			if err != nil {
				return nil, err
			}
		}

		resp, err = s.pds.Authenticate(ctx, &service.AuthenticationParams{
			Password: "",
			User:     user,
			Device:   l.Device,
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, err
	}
	return resp, nil
}

func (s *PassportService) Login(ctx context.Context, l *LoginAction) (*AuthenticationResult, error) {
	var (
		resp *service.AuthenticationInfo
		err  error
	)

	switch l.Type {
	case LocalPasswordLoginType:
		var user *userentity.User
		for _, searcher := range s.userSearchers {
			user, err = searcher.GetUser(ctx, l.Account)
			if err != nil && !errors.IsNotFound(err) {
				return nil, err
			} else if err == nil {
				break
			}
		}
		if user == nil {
			return nil, errors.NotFound(reason.UserNotExist)
		}
		resp, err = s.pds.Authenticate(ctx, &service.AuthenticationParams{
			Password: l.Password,
			User:     user,
			Device:   l.Device,
		})
	case OAuth2LoginType:
		resp, err = s.oauth2Login(ctx, l)
	default:
		return nil, errors.BadRequest(reason.LoginTypeNotSupport)
	}
	if err != nil {
		return nil, err
	}

	var (
		refreshToken string
		accessToken  string
	)
	err = s.tx.ExecTx(ctx, func(ctx context.Context) error {
		err = s.pdr.Create(ctx, resp.User.UserID, l.Device)
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
	return &AuthenticationResult{
		AccessToken:    accessToken,
		AccessTokenTTL: int(resp.AccessToken.ExpiresAt.Sub(resp.AccessToken.IssuedAt.Time).Seconds()),
		RefreshToken:   refreshToken,
		User:           resp.User,
	}, nil
}

func (s *PassportService) Logout(ctx context.Context, userInfo *entity.UserInfo) error {
	return s.deleteDeviceCore(ctx, userInfo.UserID, userInfo.DeviceID)
}

func (s *PassportService) Refresh(ctx context.Context, refreshToken string) (*AuthenticationResult, error) {
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
	return &AuthenticationResult{
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

func (s *PassportService) Init(ctx context.Context) error {
	hashedPassword, err := s.pds.GeneratePassword(defaultPassword)
	if err != nil {
		return err
	}
	_, err = s.ur.GetByEmail(ctx, defaultEmail)
	switch {
	case err == nil:
		return nil
	case errors.IsNotFound(err):
		return s.ur.Create(ctx, &userentity.User{
			UserID:   1,
			Name:     defaultUserName,
			Password: hashedPassword,
			Email:    defaultEmail,
			Mobile:   defaultMobile,
			Account:  defaultAccountName,
			ActiveInfo: &userentity.ActiveInfo{
				LastLoginAt: time.Now(),
				IP:          "127.0.0.1",
			},
		})
	default:
		return err
	}
}
