package service

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
)

type TimeoutFailOverSMSService struct {
	svcs              []SMSService
	svcIdx            int64
	timeoutCountLimit int64
	timeoutCount      int64
}

var _ SMSService = &TimeoutFailOverSMSService{}

// NewTimeoutFailOverSMSService
// 参数timeoutCountLimit表示所有goroutine(s)超时次数的限制
func NewTimeoutFailOverSMSService(svcs []SMSService, timeoutCountLimit int64) *TimeoutFailOverSMSService {
	return &TimeoutFailOverSMSService{
		svcs:              svcs,
		timeoutCountLimit: timeoutCountLimit,
	}
}

// Send goroutine(s)发送短信，goroutine(s)超时的次数超过限制后，切换到下一个短信服务商。
// 假设超时限制次数是3，此时有3个goroutine在调用这个函数，如果3个goroutine 3次调用均超时，
// 那么对于这3个goroutine，它们的第4次调用将会是下一个短信服务商；对于此时新来的goroutine，
// 它的第1次调用也会是下一个短信服务商。
// 假设超时限制次数是3，此时只有1个goroutine在调用这个函数，此goroutine调用了3次均超时，
// 那么对于这个goroutine，它的第4次调用将会是下一个短信服务商；对于此时新来的goroutine，
// 它的第1次调用也会是下一个短信服务商。
func (t *TimeoutFailOverSMSService) Send(ctx context.Context, templateID string, args []string, phoneNumbers ...string) (httpCode int, err error) {
	timeoutCount := atomic.LoadInt64(&t.timeoutCount)
	if timeoutCount >= t.timeoutCountLimit {
		atomic.CompareAndSwapInt64(&t.svcIdx, t.svcIdx, (t.svcIdx+1)%int64(len(t.svcs)))
		atomic.CompareAndSwapInt64(&t.timeoutCount, timeoutCount, 0)
	}
	httpCode, err = t.svcs[t.svcIdx].Send(ctx, templateID, args, phoneNumbers...)
	switch {
	case err == nil:
		return http.StatusOK, nil
	case errors.Is(err, context.DeadlineExceeded):
		atomic.AddInt64(&t.timeoutCount, 1)
	}

	return httpCode, err
}
