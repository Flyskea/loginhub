package messagetemplate

import (
	"context"

	"loginhub/internal/infra/data"
	messagetemplate "loginhub/internal/service/message_template"
)

var _ messagetemplate.MessageTemplateRepository = (*MessageTemplateRepo)(nil)

const (
	registerTemplateSubject = "注册验证码"
	registerTemplateBody    = `尊敬的用户，您正在注册账号。验证码为：{{.Code}}，有效期为30分钟。`
)

type MessageTemplateRepo struct {
	txm *data.TXManager
}

func NewMessageTemplateRepo(
	txm *data.TXManager,
) *MessageTemplateRepo {
	return &MessageTemplateRepo{
		txm: txm,
	}
}

func (r *MessageTemplateRepo) RegisterTemplate(
	ctx context.Context,
	messageType messagetemplate.MessageType,
) (*messagetemplate.Message, error) {
	return &messagetemplate.Message{
		Subject: registerTemplateSubject,
		Body:    registerTemplateBody,
	}, nil
}
