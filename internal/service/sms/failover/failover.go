package failover

import (
	"context"
	"errors"
	"main/internal/service/sms"
	"net/http"
	"sync"
	"sync/atomic"
)

type Service struct {
	svcs        []sms.Service
	validSVCIdx int64
	mu          sync.Mutex
}

var _ sms.Service = &Service{}

func NewService(svcs ...sms.Service) *Service {
	return &Service{svcs: svcs, validSVCIdx: 0}
}

func (s *Service) Send(ctx context.Context, templateID string,
	args []string, phoneNumbers ...string) (httpCode int, err error) {

	// 这段代码有bug，但是影响不大，写出无bug的实现，可读性不好，注释要写一大堆，不好维护，见下面的注释
	for retryCount := 0; retryCount < len(s.svcs); retryCount++ {
		_, err := s.svcs[atomic.LoadInt64(&s.validSVCIdx)].Send(ctx, templateID, args, phoneNumbers...)
		if err == nil {
			return http.StatusOK, nil
		}

		// int64很大，不用担心溢出
		atomic.StoreInt64(&s.validSVCIdx, (s.validSVCIdx+1)%int64(len(s.svcs)))
	}
	return http.StatusInternalServerError, errors.New("所有短信服务商都不可用")

	// todo: atmoic并不会阻塞，下面的解释是否有误？
	// 0  	for retryCount := 0; retryCount < len(s.svcs); retryCount++ {
	// 1		_, err := s.svcs[atomic.LoadInt64(&s.validSVCIdx)].Send(ctx, templateID, args, phoneNumbers...)
	// 2		if err == nil {
	// 3			return http.StatusOK, nil
	// 4		}
	// 5		atomic.StoreInt64(&s.validSVCIdx, (s.validSVCIdx+1)%int64(len(s.svcs)))
	// 6	}
	//
	// 假设应用刚刚启动，还没有goroutine执行过这个函数，此时s.validSVCIdx为0
	// 现在有两个goroutine同时执行上面的代码，称为A和B，他们排队拿到锁，都执行完第1句，有err，然后同时走到第5句，
	// A拿到了锁，B阻塞，A执行完第5句，释放锁，此时s.validSVCIdx为1，
	// 紧接着B拿到锁，A进入下一轮循环，正想走第1句，但此时B还没执行完第5句，锁没释放，
	// A阻塞，然后B执行完第5句，释放锁，此时s.validSVCIdx为2，A拿到锁，去读s.validSVCIdx，问题来了，本应该是去执行
	// svcs[1].Send，结果实际是执行svcs[2].Send。不过影响也不大，发送给其他服务商也是可以的。

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
