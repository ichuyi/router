package common

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
)

type HandlerFunc func(ctx *Context)

type Router struct {
	Children    map[string]*Router
	Handler     []HandlerFunc
	RequestType reflect.Type
}

func NewRouter() *Router {
	return &Router{
		Children: map[string]*Router{},
		Handler:  []HandlerFunc{},
	}
}

func (r *Router) AddChild(path string, handlerFunc HandlerFunc) *Router {
	if path == "" || !strings.HasPrefix(path, "/") {
		return nil
	}
	paths := strings.Split(path, "/")
	currentRouter := r
	for i := 1; i < len(paths); i++ {
		currentRouter.Children[paths[i]] = NewRouter()
		currentRouter = currentRouter.Children[paths[i]]
	}
	currentRouter.Handler = append(currentRouter.Handler, handlerFunc)
	return currentRouter
}
func (r *Router) RegisterType(data interface{}) {
	r.RequestType = reflect.ValueOf(data).Type()
}
func (r *Router) getRouter(path string) *Router {
	if !strings.HasPrefix(path, "/") {
		return nil
	}
	paths := strings.Split(path, "/")
	currentRouter := r
	ok := false
	for i := 1; i < len(paths); i++ {
		currentRouter, ok = currentRouter.Children[paths[i]]
		if !ok {
			return nil
		}
	}
	return currentRouter
}
func (r *Router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	path := request.RequestURI
	router := r.getRouter(path)
	if router == nil {
		return
	}
	//@todo 解析请求数据
	data := reflect.New(router.RequestType).Interface()
	request.ParseForm()
	_ = json.NewDecoder(request.Body).Decode(data)
	context := NewContext(response, request, data)
	for _, h := range router.Handler {
		h(context)
	}
	context.WriteToResponse()
	contextPool.Put(context)
}
