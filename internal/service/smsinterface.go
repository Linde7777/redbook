package service

import "context"

type SMSService interface {
	Send(ctx context.Context, templateID string, args []string, phoneNumbers ...string) (httpCode int, err error)
}
