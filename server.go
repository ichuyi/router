package main

import (
	"learn/common/LogUtil"
	"learn/common/lgo"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func hello(ctx *lgo.Context) {
	ctx.Result = Response{
		Code:    0,
		Message: "",
		Data:    "hello",
	}
}
func hi(ctx *lgo.Context) {
	ctx.Result = Response{
		Code:    0,
		Message: "",
		Data:    "hi",
	}
}
func hhh(ctx *lgo.Context) {
	ctx.Result = Response{
		Code:    0,
		Message: "",
		Data:    "hhh",
	}
}
func main() {
	l := LogUtil.NewLogger("server")
	lgo.SetLog(l)
	router := lgo.NewRouter()
	router.AddRoute("/common/hello").AddHandler(hello).RegisterType(new(Response))
	router.AddRoute("/hhh/hi").AddHandler(hi).RegisterType(new(Response))
	router.AddRoute("/hi/hhh").AddHandler(hhh).RegisterType(new(Response))
	r := router.AddRoute("/common1")
	r.AddRoute("/ddd").AddHandler(hello)
	r.AddRoute("/hi").AddHandler(hi)
	router.PrintTree("")
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	server.ListenAndServe()
}
