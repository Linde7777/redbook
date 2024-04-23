package sms

import "context"

type Service interface {
	Send(ctx context.Context, templateID string, args []string, phoneNumbers ...string) (httpCode int, err error)
}