package main

import (
	"learn/common"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func hello(ctx *common.Context) {
	ctx.Result = Response{
		Code:    0,
		Message: "",
		Data:    "hello",
	}
}
func hi(ctx *common.Context) {
	ctx.Result = Response{
		Code:    0,
		Message: "",
		Data:    "hi",
	}
}
func hhh(ctx *common.Context) {
	ctx.Result = Response{
		Code:    0,
		Message: "",
		Data:    "hhh",
	}
}
func main() {
	router := common.NewRouter()
	router.AddChild("/common/hello", hello).RegisterType(new(Response))
	router.AddChild("/hhh/hi", hi).RegisterType(new(Response))
	router.AddChild("/hi/hhh", hhh).RegisterType(new(Response))
	server := http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	server.ListenAndServe()
}
