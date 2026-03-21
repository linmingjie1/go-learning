package demo

import (
	"encoding/json"
	"fmt"
	"go-learning/middlewares"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func RunHttpRequest() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
		Handler: middlewares.Chain(
			http.DefaultServeMux,
			middlewares.WithAuth(),
			middlewares.WithTimeout(),
		),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Fprintf 会把格式化后的字符串写入到 w（响应输出流）中。
		// 这里 w 是 http.ResponseWriter，所以浏览器收到的响应体就是 "hello world"。
		_, err := fmt.Fprintf(w, "hello world")
		if err != nil {
			return
		}
	})

	http.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		id := query["id"]
		log.Println(id)
		name := query.Get("name")
		log.Println(name)
	})

	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		// ContentLength 是客户端声明的请求体长度（字节数），
		// 这里用它来提前分配接收请求体的切片。
		length := r.ContentLength
		body := make([]byte, length)

		// 从请求体流 r.Body 中读取数据到 body。
		// 读取完成后，body 中就是客户端提交的原始内容。
		_, err := r.Body.Read(body)

		// 将收到的请求体原样写回响应中，便于验证 /post 是否成功接收数据。
		_, err1 := fmt.Print(string(body))
		if err != nil || err1 != nil {
			return
		}
	})

	/*
		curl -X POST "http://127.0.0.1:8080/post/form" -d "name=Alice&age=18"
	*/
	http.HandleFunc("/post/form", func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err == nil {
			name := request.PostForm.Get("name")
			log.Println(name)
		}
	})

	/*
		curl -X POST "http://127.0.0.1:8080/multipart/form" \
		    -F "name=yibhou" \
		    -F "age=18"
	*/
	http.HandleFunc("/multipart/form", func(writer http.ResponseWriter, request *http.Request) {
		//err := request.ParseMultipartForm(32 << 20)
		//if err == nil {
		//	log.Println(request.MultipartForm)
		//}

		log.Println(request.FormValue("name"))
	})

	// 上传文件
	/*
		curl -X POST "http://127.0.0.1:8080/upload" -F "file=@README.md"
	*/
	http.HandleFunc("/upload", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(writer, "only POST is supported", http.StatusMethodNotAllowed)
			return
		}

		err := request.ParseMultipartForm(10 << 20) // 最多解析 10MB 到内存，超出的部分会落盘临时文件
		if err != nil {
			http.Error(writer, "invalid multipart form", http.StatusBadRequest)
			return
		}
		defer func() {
			// 清理 ParseMultipartForm 产生的临时文件
			if err := request.MultipartForm.RemoveAll(); err != nil {
				log.Printf("cleanup multipart temp files failed: %v", err)
			}
		}()

		file, header, err := request.FormFile("file")
		if err != nil {
			http.Error(writer, "file field is required", http.StatusBadRequest)
			return
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Printf("close uploaded file failed: %v", err)
			}
		}()

		// 递归创建 uploads 目录；目录已存在时不会报错，0o755(八进制)，表示owner可读写执行，group或others可读执行
		err = os.MkdirAll("uploads", 0o755)
		if err != nil {
			http.Error(writer, "failed to create upload dir", http.StatusInternalServerError)
			return
		}
		// 创建目标文件
		filename := filepath.Base(header.Filename)
		savePath := filepath.Join("uploads", filename)
		dst, err := os.Create(savePath)
		if err != nil {
			http.Error(writer, "failed to create destination file", http.StatusInternalServerError)
			return
		}
		defer func() {
			if err := dst.Close(); err != nil {
				log.Printf("close destination file failed: %v", err)
			}
		}()

		// 将读取到的数据写入目标文件
		n, err := io.Copy(dst, file)
		if err != nil {
			http.Error(writer, "failed to save file", http.StatusInternalServerError)
			return
		}

		log.Printf("upload success: %s (%d bytes)", filename, n)
		_, _ = fmt.Fprintf(writer, "upload success: %s (%d bytes)\n", filename, n)
	})

	type Student struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	// json
	http.HandleFunc("/json", func(writer http.ResponseWriter, request *http.Request) {
		stu := Student{}

		switch request.Method {
		case http.MethodPost:
			decoder := json.NewDecoder(request.Body)
			err := decoder.Decode(&stu) // 从请求体中解析 JSON 数据到 stu 结构体
			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
				return
			}

			writer.Header().Set("Content-Type", "application/json")
			encoder := json.NewEncoder(writer)
			err = encoder.Encode(stu) // 将 stu 结构体编码为 JSON 格式写入响应体
			if err != nil {
				writer.WriteHeader(http.StatusBadRequest)
			}
		default:
			http.Error(writer, "only POST is supported", http.StatusMethodNotAllowed)
		}
	})

	// 测试超时
	http.HandleFunc("/timeout", func(writer http.ResponseWriter, request *http.Request) {
		// 模拟一个耗时操作
		time.Sleep(10 * time.Second)
	})

	err := server.ListenAndServe()
	if err != nil {
		return
	}
}
