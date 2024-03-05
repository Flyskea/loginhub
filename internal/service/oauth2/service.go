package oauth2

import (
	"context"
	"net/http"

	"github.com/Flyskea/gotools/errors"

	"loginhub/internal/base/pager"
	"loginhub/internal/base/reason"
	"loginhub/internal/domain/oauth2/entity"
	"loginhub/internal/domain/oauth2/repository"
	"loginhub/pkg/random"
)

const (
	tmpSessionLength = 32
)

type OAuth2Service struct {
	opir repository.OAuth2ProviderInfoRepository
}

func NewOauth2Service(
	opir repository.OAuth2ProviderInfoRepository,
) *OAuth2Service {
	return &OAuth2Service{
		opir: opir,
	}
}

func (s *OAuth2Service) CreateProviderInfo(ctx context.Context, provider *entity.ProviderInfo) error {
	_, err := s.opir.GetByType(ctx, provider.Type)
	if err == nil {
		return errors.BadRequest(reason.OAuth2ProviderDuplicate)
	}
	if !entity.IsSupportedProvider(provider.Type) {
		return errors.BadRequest(reason.OAuth2ProviderNotImplemented)
	}
	return s.opir.Save(ctx, provider)
}

func (s *OAuth2Service) ListProviderInfos(ctx context.Context, pageCond *pager.PageCond) ([]*entity.ProviderInfo, int64, error) {
	return s.opir.List(ctx, pageCond)
}

func (s *OAuth2Service) GetProviderInfo(ctx context.Context, providerType string) (*entity.ProviderInfo, error) {
	info, err := s.opir.GetByType(ctx, providerType)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NotFound(reason.OAuth2ProviderNotExist)
		}
		return nil, err
	}
	return info, nil
}

func (s *OAuth2Service) SaveProviderInfo(ctx context.Context, provider *entity.ProviderInfo) error {
	return s.opir.Save(ctx, provider)
}

func (s *OAuth2Service) UpdateProviderInfo(ctx context.Context, provider *entity.ProviderInfo) error {
	_, err := s.opir.GetByType(ctx, provider.Type)
	if err != nil {
		if errors.IsNotFound(err) {
			return errors.BadRequest(reason.OAuth2ProviderNotExist)
		}
		return err
	}
	return s.opir.Update(ctx, provider)
}

func (s *OAuth2Service) DeleteProviderInfo(ctx context.Context, id int64) error {
	return s.opir.DeleteByID(ctx, id)
}

// OAuth2RequestURL build oauth2 request url
// return request url and random tmp session
func (s *OAuth2Service) OAuth2RequestURL(ctx context.Context, providerType string) (string, string, error) {
	providerInfo, err := s.GetProviderInfo(ctx, providerType)
	if err != nil {
		return "", "", err
	}
	provider, err := entity.NewOauth2Provider(
		providerInfo.Type,
		providerInfo.ClientID,
		providerInfo.ClientSecret,
		providerInfo.RedirectURL,
	)
	if err != nil {
		return "", "", err
	}
	requestURL, state, err := provider.BuildRequestURL()
	if err != nil {
		return "", "", err
	}

	sessionID := random.RandomString(random.Alphanumeric, tmpSessionLength)
	err = s.opir.SaveState(ctx, providerType, sessionID, state)
	if err != nil {
		return "", "", err
	}
	return requestURL, sessionID, nil
}

func (s *OAuth2Service) OAuth2UserInfo(ctx context.Context, action *GetUserInfoAction) (*GetUserInfoResult, error) {
	state, err := s.opir.GetStateByKey(ctx, action.ProviderType, action.OAuth2SessionID)
	switch {
	case err == nil:
		if state != action.RequestState {
			return nil, errors.BadRequest(reason.OAuth2StateMismatch)
		}
	case errors.IsNotFound(err):
		return nil, errors.BadRequest(reason.OAuth2StateNotFound)
	default:
		return nil, err
	}

	providerInfo, err := s.GetProviderInfo(ctx, action.ProviderType)
	if err != nil {
		return nil, err
	}
	provider, err := entity.NewOauth2Provider(
		providerInfo.Type,
		providerInfo.ClientID,
		providerInfo.ClientSecret,
		providerInfo.RedirectURL,
	)
	if err != nil {
		return nil, err
	}
	token, err := provider.GetToken(ctx, action.Code)
	if err != nil {
		return nil, err
	}
	provider.SetClient(http.DefaultClient)
	userInfo, err := provider.GetUserInfo(ctx, token)
	if err != nil {
		return nil, err
	}

	return &GetUserInfoResult{
		Token:    token,
		UserInfo: userInfo,
	}, nil
}
