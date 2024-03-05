package email

import (
	"context"
	"fmt"
	"log/slog"
	"net/smtp"
	"time"

	"github.com/Flyskea/gotools/errors"
	"github.com/jordan-wright/email"

	"loginhub/internal/base/reason"
	"loginhub/internal/conf"
	"loginhub/pkg/convert"
)

const (
	timeout = 10 * time.Second
)

type EmailSender struct {
	c      *email.Pool
	logger *slog.Logger
	conf   *conf.SMTP
}

func NewEmailSender(
	c *conf.SMTP,
	logger *slog.Logger,
) (*EmailSender, func(), error) {
	client, err := email.NewPool(
		fmt.Sprintf("%s:%d", c.GetHost(), c.GetPort()),
		5,
		smtp.PlainAuth("", c.GetUserName(), c.GetPassword(), c.GetHost()),
	)
	if err != nil {
		return nil, nil, errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	s := &EmailSender{
		c:      client,
		conf:   c,
		logger: logger,
	}

	cleanup := func() {
		client.Close()
	}

	return s, cleanup, nil
}

func (s *EmailSender) Send(ctx context.Context, toEmailAddr, subject, body string) error {
	m := email.NewEmail()
	m.From = fmt.Sprintf("%s <%s>", s.conf.GetFromName(), s.conf.GetUserName())
	m.To = append(m.To, toEmailAddr)
	m.Subject = subject
	m.Text = convert.StringToBytes(body)
	m.Headers.Add("Content-Type", "text/plain; charset=UTF-8")
	err := s.c.Send(m, timeout)
	if err != nil {
		return errors.InternalServer(reason.UnknownError).WithError(err).WithStack()
	}
	return nil
}
