package captcha

import (
	"bytes"
	"context"
	"math/rand"
	"time"

	"github.com/Flyskea/gotools/errors"

	"loginhub/internal/base/reason"
	"loginhub/internal/infra/email"
	messagetemplate "loginhub/internal/service/message_template"
)

type CaptchaType uint8

const (
	RegisterCaptchatType CaptchaType = iota + 1
)

const (
	codeMap = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type CaptchaRepository interface {
	Save(ctx context.Context, userID int64, capchaType CaptchaType, c *Captcha) error
	SaveByEmail(ctx context.Context, email string, capchaType CaptchaType, c *Captcha) error
	Get(ctx context.Context, userID int64, capchaType CaptchaType) (*Captcha, error)
	GetByEmail(ctx context.Context, email string, capchaType CaptchaType) (*Captcha, error)
}

type CaptchaService struct {
	cr          CaptchaRepository
	emailSender *email.EmailSender
	mts         *messagetemplate.MessageTemplateService
}

func New(
	cr CaptchaRepository,
	emailSender *email.EmailSender,
	mts *messagetemplate.MessageTemplateService,
) *CaptchaService {
	return &CaptchaService{
		cr:          cr,
		emailSender: emailSender,
		mts:         mts,
	}
}

func (s *CaptchaService) genCaptcha() string {
	rand.New(rand.NewSource(time.Now().Unix()))
	code := bytes.Buffer{}
	code.Grow(6)
	codeMapLen := len(codeMap)
	for i := 0; i < 6; i++ {
		code.WriteByte(codeMap[rand.Intn(codeMapLen)])
	}
	return code.String()
}

func (s *CaptchaService) isCodeCorrect(
	code string,
	capcha *Captcha,
	err error,
) error {
	switch {
	case err != nil && errors.IsNotFound(err):
		return errors.BadRequest(reason.CaptchaNotExist)
	case err != nil:
		return err
	case capcha.IsCorrect(code):
		return nil
	default:
		return errors.BadRequest(reason.CaptchaMismatch)
	}
}

func (s *CaptchaService) IsCodeCorrect(
	ctx context.Context,
	userID int64,
	capchaType CaptchaType,
	code string,
) error {
	capcha, err := s.cr.Get(ctx, userID, capchaType)
	return s.isCodeCorrect(code, capcha, err)
}

func (s *CaptchaService) IsCodeCorrectByEmail(
	ctx context.Context,
	email string,
	capchaType CaptchaType,
	code string,
) error {
	capcha, err := s.cr.GetByEmail(ctx, email, capchaType)
	return s.isCodeCorrect(code, capcha, err)
}

func (s *CaptchaService) EmailCaptcha(ctx context.Context, email string, capchaType CaptchaType) (string, error) {
	captcha, err := s.cr.GetByEmail(ctx, email, capchaType)
	if err != nil && !errors.IsNotFound(err) {
		return "", err
	}
	if captcha != nil && !captcha.Allow() {
		return "", errors.BadRequest(reason.CaptchaReachSendLimit)
	}
	code := s.genCaptcha()
	err = s.cr.SaveByEmail(ctx, email, capchaType, NewCaptcha(code))
	if err != nil {
		return "", err
	}

	var msg *messagetemplate.Message
	switch capchaType {
	case RegisterCaptchatType:
		msg, err = s.mts.RegisterTemplate(ctx, code, messagetemplate.Email)
	default:
		return "", errors.BadRequest(reason.CaptchaTypeNotSupport)
	}
	if err != nil {
		return "", err
	}
	err = s.emailSender.Send(ctx, email, msg.Subject, msg.Body)
	if err != nil {
		return "", err
	}

	return code, nil
}
