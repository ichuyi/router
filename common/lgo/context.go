package lgo

import (
	"encoding/json"
	"fmt"
	"learn/common/LogUtil"
	"net/http"
	"sync"
)

type Log interface {
	Error(s string)
	Info(s string)
	Fatal(s string)
}

var logger Log = LogUtil.StdLogger

func SetLog(log Log) {
	logger = log
}

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Input          interface{}
	Result         interface{}
}

var contextPool = &sync.Pool{New: func() interface{} {
	return &Context{
		Request:        nil,
		ResponseWriter: nil,
		Input:          nil,
	}
}}

func NewContext(response http.ResponseWriter, request *http.Request) *Context {
	c := contextPool.Get()
	context := c.(*Context)
	context.ResponseWriter = response
	context.Request = request
	return context
}
func (ctx *Context) WriteToResponse() {
	data, err := json.Marshal(ctx.Result)
	if err != nil {
		return
	}
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err = ctx.ResponseWriter.Write(data)
	if err != nil {
		logger.Error(fmt.Sprintf("write to response error: %s", err.Error()))
		return
	}
}
func (ctx *Context) NotFound() {
	ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
}
func (ctx *Context) BadRequest() {
	ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
}
