package service

import (
	"context"

	"github.com/Flyskea/gotools/errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"loginhub/internal/base/reason"
	passportentity "loginhub/internal/domain/passport/entity"
	"loginhub/internal/domain/user/repository"
	"loginhub/pkg/convert"
)

type PassportService struct {
	ur repository.UserRepository
}

func NewPassportService(
	ur repository.UserRepository,
) *PassportService {
	return &PassportService{
		ur: ur,
	}
}

func (s *PassportService) Authentication(ctx context.Context, info *LoginInfo) (*AuthenticationInfo, error) {
	info.Device.ID = convert.TrimUUID(uuid.NewString())
	user, err := s.ur.GetByEmail(ctx, info.Email)
	if err != nil && !errors.IsBadRequest(err) {
		return nil, err
	}
	if err != nil {
		user, err = s.ur.GetByMobile(ctx, info.Phone)
		if err != nil && !errors.IsBadRequest(err) {
			return nil, err
		}
	}
	if err != nil {
		return nil, errors.BadRequest(reason.AccountOrPasswordWrong)
	}
	ok, err := s.ComparePassword(info.Password, user.Password)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.BadRequest(reason.AccountOrPasswordWrong)
	}

	accessToken := passportentity.CreateAccessToken(user.UserID, user.Name, info.Device.ID)
	refreshToken := passportentity.CreateRefreshToken(user.UserID, user.Name, info.Device.ID)
	return &AuthenticationInfo{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *PassportService) Register(ctx context.Context, info *RegisterInfo) (*AuthenticationInfo, error) {
	err := validateUserName(info.User.Name)
	if err != nil {
		return nil, err
	}
	err = validatePassword(info.User.Password)
	if err != nil {
		return nil, err
	}
	hashed, err := s.GeneratePassword(info.User.Password)
	if err != nil {
		return nil, err
	}

	user := info.User
	user.Password = hashed
	switch info.RegisterType {
	case EmailRegisterType:
		err = validateEmail(user.Email)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.BadRequest(reason.RegisterTypeNotSupport)
	}

	accessToken := passportentity.CreateAccessToken(user.UserID, user.Name, info.Device.ID)
	refreshToken := passportentity.CreateRefreshToken(user.UserID, user.Name, info.Device.ID)
	return &AuthenticationInfo{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *PassportService) Refresh(ctx context.Context, refreshToken *passportentity.RefreshToken) (*AuthenticationInfo, error) {
	newRefreshToken := refreshToken.Refresh()
	user, err := s.ur.GetByUserID(ctx, refreshToken.UserID)
	if err != nil {
		return nil, err
	}
	accessToken := passportentity.CreateAccessToken(user.UserID, user.Name, newRefreshToken.DeviceID)
	return &AuthenticationInfo{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         user,
	}, nil
}

func (s *PassportService) GeneratePassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	return convert.BytesToString(hash), nil
}

func (s *PassportService) ComparePassword(password string, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(convert.StringToBytes((hashedPassword)), convert.StringToBytes(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, errors.BadRequest(reason.AccountOrPasswordWrong)
		}
		return false, errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	return true, nil
}
