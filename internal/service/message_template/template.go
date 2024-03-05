package messagetemplate

import (
	"bytes"
	"context"
	"html/template"

	"github.com/Flyskea/gotools/errors"

	"loginhub/internal/base/reason"
)

type Message struct {
	Subject string
	Body    string
}

type MessageType uint8

const (
	Email MessageType = iota + 1
)

type MessageTemplateRepository interface {
	RegisterTemplate(ctx context.Context, messageType MessageType) (*Message, error)
}

type MessageTemplateService struct {
	mtr MessageTemplateRepository
}

func NewMessageTemplateService(
	mtr MessageTemplateRepository,
) *MessageTemplateService {
	return &MessageTemplateService{
		mtr: mtr,
	}
}

func (s *MessageTemplateService) RegisterTemplate(
	ctx context.Context,
	code string,
	messageType MessageType,
) (*Message, error) {
	msg, err := s.mtr.RegisterTemplate(ctx, messageType)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(nil)
	err = template.Must(template.New("body").Parse(msg.Body)).
		Execute(body, map[string]string{"Code": code})
	if err != nil {
		return nil, errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	msg.Body = body.String()
	return msg, nil
}
