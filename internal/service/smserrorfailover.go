package service

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"sync/atomic"
)

type FailOverService struct {
	svcs        []SMSService
	validSVCIdx int64
	mu          sync.Mutex
}

var _ SMSService = &FailOverService{}

func NewFailOverService(svcs ...SMSService) *FailOverService {
	return &FailOverService{svcs: svcs, validSVCIdx: 0}
}

func (s *FailOverService) Send(ctx context.Context, templateID string,
	args []string, phoneNumbers ...string) (httpCode int, err error) {

	for retryCount := 0; retryCount < len(s.svcs); retryCount++ {
		_, err := s.svcs[atomic.LoadInt64(&s.validSVCIdx)].Send(ctx, templateID, args, phoneNumbers...)
		if err == nil {
			return http.StatusOK, nil
		}

		// int64很大，不用担心溢出
		atomic.StoreInt64(&s.validSVCIdx, (s.validSVCIdx+1)%int64(len(s.svcs)))
	}
	return http.StatusInternalServerError, errors.New("所有短信服务商都不可用")

	// 0	for retryCount := 0; retryCount < len(s.svcs); retryCount++ {
	// 1		tempValidSVCIdx := s.validSVCIdx
	// 2		s.mu.Unlock()
	// 3		_, err := s.svcs[tempValidSVCIdx].Send(ctx, templateID, args, phoneNumbers...)
	// 4		if err == nil {
	// 5			return http.StatusOK, nil
	// 6		}
	// 7
	// 8		beforeLockIdx := s.validSVCIdx
	// 9		s.mu.Lock()
	// 10		if beforeLockIdx == s.validSVCIdx {
	// 11			s.validSVCIdx = (s.validSVCIdx + 1) % int64(len(s.svcs))
	// 12		}
	// 15	}
	//
	// 假设应用刚启动，还没有goroutine执行过这段代码，s.validSVCIdx为0，
	// 现在有两个goroutine A和B同时进入第0轮循环
	// （第1句和第2句的作用在第1轮循环后会体现），A和B执行完第3句，err不为空，
	// 执行完第8句，A和B的beforeLockIdx均为0，A和B都将要执行第9句，
	// A拿到了锁，此时s.validSVCIdx为0，和A的beforeIdx相同，所以会进入第11句，
	// 执行完后s.validSVCIdx为1。10到12句的这个判断到底有什么用，后面会体现。
	// A进入第1轮循环，记录一下s.ValidSVCIdx后释放锁，执行Send。
	// 第1句和第2句就是让A能够尽快释放锁，而不是等到Send执行完后才释放，
	// A执行完第2句，释放锁。B现在拿到锁，进入第10句，B的beforeLockIdx是0，
	// 而s.ValidSVCIdx是1，不相等，所以不用再让s.ValidSVCIdx自增。

}
