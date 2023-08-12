package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	"sync"
	"time"
)

type errResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}
type buffpool struct {
	pool sync.Pool
}

func (p *buffpool) Get() *bytes.Buffer {
	buf := p.pool.Get()
	if buf == nil {
		return &bytes.Buffer{}
	}
	return buf.(*bytes.Buffer)
}

// Put a bytes.Buffer pointer to BufferPool
func (p *buffpool) Put(buf *bytes.Buffer) {
	p.pool.Put(buf)
}

type TimeoutWriter struct {
	gin.ResponseWriter
	// body
	body *bytes.Buffer
	// header
	// 变动点2:  让子协程和父协程分别写不同的header
	h http.Header

	mu          sync.Mutex
	timedOut    bool
	wroteHeader bool
	// 变动点3:  让子协程和父协程分别写不同的code
	code int
}

func (tw *TimeoutWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		//return 0, http.ErrHandlerTimeout
		// 已经超时了，就不再写数据
		return 0, nil
	}

	return tw.body.Write(b)
}

func (tw *TimeoutWriter) WriteHeader(code int) {
	fmt.Println("----xxx---", "TimeoutWriter-WriteHeader")
	checkWriteHeaderCode(code)
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return
	}
	tw.writeHeader(code)
}

func (tw *TimeoutWriter) writeHeader(code int) {
	tw.wroteHeader = true
	tw.code = code
}

func (tw *TimeoutWriter) WriteHeaderNow() {
	fmt.Println("----xxx---", "TimeoutWriter-WriteHeaderNow")
}

func (tw *TimeoutWriter) Header() http.Header {
	return tw.h
}

func checkWriteHeaderCode(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", code))
	}
}

func Timeout(t time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// wrap the request context with a timeout
		// sync.Pool
		b := buffpool{}
		buffer := b.Get()

		tw := &TimeoutWriter{body: buffer, ResponseWriter: c.Writer, h: make(http.Header)}
		c.Writer = tw

		ctx, cancel := context.WithTimeout(c.Request.Context(), t)
		c.Request = c.Request.WithContext(ctx)

		// channel 容量必须大于0
		// 否则母协程因超时退出，子协程可能永远无法退出
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)
		go func() {
			// 变动点1: 增加子协程的recover
			defer func() {
				if p := recover(); p != nil {
					fmt.Println("handler error", p)
					panicChan <- p
				}
			}()

			c.Next()
			finish <- struct{}{}
		}()

		select {
		case p := <-panicChan:
			panic(p)
		case <-ctx.Done():
			tw.mu.Lock()
			defer tw.mu.Unlock()

			tw.ResponseWriter.WriteHeader(http.StatusServiceUnavailable)
			bt, _ := json.Marshal(errResponse{Code: "503",
				Msg: http.ErrHandlerTimeout.Error()})
			tw.ResponseWriter.Write(bt)
			c.Abort()
			cancel()
			tw.timedOut = true
			// 如果超时的话，buffer无法主动清除，只能等待GC回收
		case <-finish:
			tw.mu.Lock()
			defer tw.mu.Unlock()
			dst := tw.ResponseWriter.Header()
			for k, vv := range tw.Header() {
				dst[k] = vv
			}
			fmt.Println("tw.code", tw.code)
			tw.ResponseWriter.WriteHeader(tw.code)
			tw.ResponseWriter.Write(buffer.Bytes())
			b.Put(buffer)
		}
	}
}

func short(c *gin.Context) {
	time.Sleep(1 * time.Second)
	// 子协程操作的header，其实是TimeoutWriter中的Header
	c.JSON(http.StatusOK, gin.H{"hello": "world"})
}

func nocontent(c *gin.Context) {
	//c.Status(204)
	time.Sleep(1 * time.Second)
	c.Data(http.StatusNoContent, "", []byte{})
}
