package middlewares

import (
	"context"
	"net/http"
	"time"
)

type TimeoutMiddleware struct {
	Next http.Handler
}

func (middleware TimeoutMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if middleware.Next == nil {
		middleware.Next = http.DefaultServeMux
	}

	// 基于当前请求上下文派生一个带 5 秒超时的新上下文，
	// 后续链路如果感知 ctx.Done()，就可以主动结束处理。
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	r = r.WithContext(ctx) // 让下游使用新的带超时的context

	// 用 channel 接收下游处理完成的信号，便于和超时信号一起等待。
	// chan struct{}表示传递类型是空的结构体
	ch := make(chan struct{})
	// 启动一个新的 goroutine，异步执行这个匿名函数。
	go func() {
		middleware.Next.ServeHTTP(w, r)
		ch <- struct{}{}
	}()
	// 主 goroutine 等两个结果之一：下游处理完成的信号，或者超时信号。
	select {
	case <-ch:
		// 下游在超时前完成，请求正常结束，主动释放定时器资源。
		return
	case <-ctx.Done():
		// 超时后直接返回 504，表示网关/中间件等待下游响应超时。
		w.WriteHeader(http.StatusGatewayTimeout)
	}
}
