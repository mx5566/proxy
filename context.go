package proxy

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
}

// 重置
func (this *Context) Reset(rw http.ResponseWriter, r *http.Request) {
	this.ResponseWriter = rw
	this.Request = r

}

// 响应字符串
func (this *Context) WriteString(content string) {
	_, _ = this.ResponseWriter.Write([]byte(content))
}

func (this *Context) Header(key, value string) {
	this.ResponseWriter.Header().Set(key, value)
}

func (this *Context) ServerJson(data interface{}) {
	var content []byte
	var err error
	content, err = json.Marshal(data)

	if err != nil {
		http.Error(this.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	this.Header("Content-Type", "application/json; charset=utf-8")
	this.Header("Content-Length", strconv.Itoa(len(content)))

	var buf = &bytes.Buffer{}
	buf.Write(content)

	_, _ = io.Copy(this.ResponseWriter, buf)

}
