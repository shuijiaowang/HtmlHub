package middleware

import (
	"htmlhub/util/response"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	loginFailWindow    = time.Minute
	loginFailThreshold = 3
)

type ipLoginFailRecord struct {
	FailCount  int
	FirstFail  time.Time
	BlockUntil time.Time
}

var (
	ipLoginFailStore = make(map[string]*ipLoginFailRecord)
	ipLoginFailMu    sync.Mutex
)

// LoginIPRateLimit 仅用于登录/注册路由，实际只统计 /api/user/login 的失败次数
func LoginIPRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只对登录接口执行限流逻辑
		if c.FullPath() != "/api/user/login" {
			c.Next()
			return
		}

		now := time.Now()
		ip := c.ClientIP()

		ipLoginFailMu.Lock()
		record := ipLoginFailStore[ip]
		if record != nil && now.Before(record.BlockUntil) {
			ipLoginFailMu.Unlock()
			response.FailWithMessage("登录失败次数过多，请1分钟后再试", c)
			c.Abort()
			return
		}
		ipLoginFailMu.Unlock()

		c.Next()

		failed, ok := c.Get("login_failed")
		if !ok || failed != true {
			// 登录成功后清理该IP失败记录
			ipLoginFailMu.Lock()
			delete(ipLoginFailStore, ip)
			ipLoginFailMu.Unlock()
			return
		}

		now = time.Now()
		ipLoginFailMu.Lock()
		record = ipLoginFailStore[ip]
		if record == nil {
			record = &ipLoginFailRecord{
				FailCount: 1,
				FirstFail: now,
			}
			ipLoginFailStore[ip] = record
			ipLoginFailMu.Unlock()
			return
		}

		if now.Sub(record.FirstFail) > loginFailWindow {
			record.FailCount = 1
			record.FirstFail = now
			record.BlockUntil = time.Time{}
			ipLoginFailMu.Unlock()
			return
		}

		record.FailCount++
		if record.FailCount > loginFailThreshold {
			record.BlockUntil = now.Add(loginFailWindow)
		}
		ipLoginFailMu.Unlock()
	}
}
