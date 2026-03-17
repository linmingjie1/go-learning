package demo

import "net/http"

// helloHandler 是一个自定义处理器（handler）。
// 只要类型实现了 ServeHTTP 方法，就可以处理 HTTP 请求。
type helloHandler struct{}

func (m *helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 向响应体写入内容，浏览器会看到 "Hello World"。
	_, err := w.Write([]byte("Hello World"))
	if err != nil {
		return
	}
}

// aboutHandler 和 helloHandler 一样，都是自定义 handler。
type aboutHandler struct{}

func (m *aboutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("About"))
	if err != nil {
		return
	}
}

func welcome(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Welcome"))
	if err != nil {
		return
	}
}

func test1() {
	// HandleFunc 是更简洁的写法：直接传入函数。
	// 底层会把这个函数转换为 http.HandlerFunc 类型，
	// 而 http.HandlerFunc 实现了 ServeHTTP，所以它本质上也是一个 handler。
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var data = []byte("hello world")
		_, err := w.Write(data)

		if err != nil {
			return
		}
	})

	// nil 同样表示使用默认路由器 DefaultServeMux。
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func test2() {
	mh := helloHandler{}
	ah := aboutHandler{}

	// 创建一个 HTTP 服务器实例，监听本机 8080 端口。
	// Handler 为 nil 表示使用默认路由器 DefaultServeMux。
	server := http.Server{
		Addr:    ":8080",
		Handler: nil,
	}
	// 把 URL 路径和对应 handler 绑定起来。
	// 访问 /hello 时会走 helloHandler，访问 /about 时会走 aboutHandler。
	http.Handle("/hello", &mh)
	http.Handle("/about", &ah)
	http.HandleFunc("/welcome", welcome)

	// 启动服务器并阻塞当前 goroutine，直到服务退出或报错。
	err := server.ListenAndServe()
	if err != nil {
		return
	}
}
