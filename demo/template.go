package demo

import (
    "log"
    "net/http"
    "text/template"
)

func RunTemplate() {
    server := http.Server{
        Addr: ":8080",
    }
    http.HandleFunc("/", process)

    err := server.ListenAndServe()
    if err != nil {
        return
    }
}

func process(w http.ResponseWriter, r *http.Request) {
    // 解析模板文件，得到一个 *template.Template 对象。
    files, _ := template.ParseFiles("templates/index.html")
    data := struct {
        Title   string
        Name    string
        IsLogin bool
    }{
        Title:   "Go Template Demo",
        Name:    "linmingjie",
        IsLogin: true,
    }
    err := files.Execute(w, data)
    if err != nil {
        log.Println(err)
    }
}

func loadTemplates() *template.Template {
	// 创建模板集合根节点，后续解析的模板都会挂载到这个集合中。
    result := template.New("myTemplates")
	// ParseGlob 会一次性加载 templates 目录下所有 html 模板；
	// Must 在解析失败时直接 panic，适合启动阶段尽早暴露模板错误。
    template.Must(result.ParseGlob("templates/*.html"))
    return result
}
