package common

import (
	"encoding/json"
	"net/http"
	"sync"
)

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	Input          interface{}
	Result         interface{}
}

var contextPool = sync.Pool{New: func() interface{} {
	return &Context{
		Request:        nil,
		ResponseWriter: nil,
		Input:          nil,
	}
}}

func NewContext(response http.ResponseWriter, request *http.Request, data interface{}) *Context {
	c := contextPool.Get()
	context := c.(*Context)
	context.ResponseWriter = response
	context.Request = request
	context.Input = data
	return context
}
func (ctx *Context) WriteToResponse() {
	data, err := json.Marshal(ctx.Result)
	if err != nil {
		return
	}
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, _ = ctx.ResponseWriter.Write(data)
}
