package middleware

import (
	"htmlhub/util/response"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	writeWindow    = time.Minute
	writeThreshold = 100
	writeBlockTime = time.Minute
)

type ipWriteRecord struct {
	Count      int
	WindowFrom time.Time
	BlockUntil time.Time
}

var (
	ipWriteStore = make(map[string]*ipWriteRecord)
	ipWriteMu    sync.Mutex
)

// HighRiskWriteRateLimit 用于高风险写接口，限制单IP高频提交
func HighRiskWriteRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		//|| c.Request.Method == http.MethodGet
		if c.Request.Method == http.MethodOptions || c.Request.Method == http.MethodHead {
			c.Next()
			return
		}

		now := time.Now()
		ip := c.ClientIP()

		ipWriteMu.Lock()
		record := ipWriteStore[ip]
		if record == nil {
			ipWriteStore[ip] = &ipWriteRecord{
				Count:      1,
				WindowFrom: now,
			}
			ipWriteMu.Unlock()
			c.Next()
			return
		}

		if now.Before(record.BlockUntil) {
			ipWriteMu.Unlock()
			response.FailWithMessage("请求过于频繁，请稍后再试", c)
			c.Abort()
			return
		}

		if now.Sub(record.WindowFrom) > writeWindow {
			record.Count = 1
			record.WindowFrom = now
			record.BlockUntil = time.Time{}
			ipWriteMu.Unlock()
			c.Next()
			return
		}

		record.Count++
		if record.Count > writeThreshold {
			record.BlockUntil = now.Add(writeBlockTime)
			ipWriteMu.Unlock()
			response.FailWithMessage("请求过于频繁，请稍后再试", c)
			c.Abort()
			return
		}
		ipWriteMu.Unlock()

		c.Next()
	}
}
