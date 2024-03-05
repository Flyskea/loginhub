package passport

import (
	"context"
	userentity "loginhub/internal/domain/user/entity"
	userrepo "loginhub/internal/domain/user/repository"
)

type LocalUserSearcher interface {
	GetUser(ctx context.Context, account string) (*userentity.User, error)
}

type EmailSearcher struct {
	ur userrepo.UserRepository
}

func NewEmailAccountSearcher(ur userrepo.UserRepository) *EmailSearcher {
	return &EmailSearcher{
		ur: ur,
	}
}

func (e *EmailSearcher) GetUser(ctx context.Context, account string) (*userentity.User, error) {
	return e.ur.GetByEmail(ctx, account)
}

type AccountSearcher struct {
	ur userrepo.UserRepository
}

func NewAccountSearcher(ur userrepo.UserRepository) *AccountSearcher {
	return &AccountSearcher{
		ur: ur,
	}
}

func (a *AccountSearcher) GetUser(ctx context.Context, account string) (*userentity.User, error) {
	return a.ur.GetByAccountName(ctx, account)
}

type MobileSearcher struct {
	ur userrepo.UserRepository
}

func NewMobileAccountSearcher(ur userrepo.UserRepository) *MobileSearcher {
	return &MobileSearcher{
		ur: ur,
	}
}

func (m *MobileSearcher) GetUser(ctx context.Context, account string) (*userentity.User, error) {
	return m.ur.GetByMobile(ctx, account)
}
