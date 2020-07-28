package lgo

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/schema"
	"net/http"
	"reflect"
	"strings"
)

var formDecoder = schema.NewDecoder()

func init() {
	formDecoder.IgnoreUnknownKeys(true)
}

type HandlerFunc func(ctx *Context)
type Router struct {
	Children    map[string]*Router
	Handler     []HandlerFunc
	RequestType reflect.Type
	IsLeaf      bool
	method      []string
}

var EmptyRouter = Router{}

func NewRouter() *Router {
	return &Router{
		Children:    map[string]*Router{},
		Handler:     []HandlerFunc{},
		RequestType: nil,
		IsLeaf:      true,
		method:      []string{},
	}
}
func (r *Router) AddRoute(path string) *Router {
	if path == "" || !strings.HasPrefix(path, "/") {
		logger.Error(fmt.Sprintf("path is error: %s", path))
		return nil
	}
	paths := strings.Split(path, "/")
	currentRouter := r
	for i := 1; i < len(paths); i++ {
		router := NewRouter()
		//将父级的handler加到子级，实现过滤器
		router.Handler = make([]HandlerFunc, len(currentRouter.Handler))
		copy(router.Handler, currentRouter.Handler)
		currentRouter.Children[paths[i]] = router
		currentRouter.IsLeaf = false
		currentRouter = router
	}
	return currentRouter
}
func (r *Router) AddHandler(handlerFunc HandlerFunc) *Router {
	r.Handler = append(r.Handler, handlerFunc)
	return r
}
func (r *Router) Method(method ...string) *Router {
	for i := range method {
		r.method = append(r.method, method[i])
	}
	return r
}
func (r *Router) RegisterType(data interface{}) {
	if r.IsLeaf {
		r.RequestType = reflect.ValueOf(data).Type().Elem()
	} else {
		logger.Error(fmt.Sprintf("current router is not a leaf node"))
	}
}
func (r *Router) getRouter(path, method string) *Router {
	if !strings.HasPrefix(path, "/") {
		logger.Error(fmt.Sprintf("path is error: %s", path))
		return nil
	}
	paths := strings.Split(path, "/")
	currentRouter := r
	ok := false
	for i := 1; i < len(paths); i++ {
		currentRouter, ok = currentRouter.Children[paths[i]]
		if !ok {
			logger.Info(fmt.Sprintf("path not exist: %s", path))
			return nil
		}
	}
	if currentRouter.IsLeaf && currentRouter.methodExist(method) {
		return currentRouter
	} else {
		return nil
	}
}
func (r *Router) methodExist(m string) bool {
	if len(r.method) == 0 {
		r.method = append(r.method, "ALL")
	}
	for i := range r.method {
		if m == r.method[i] || r.method[i] == "ALL" {
			return true
		}
	}
	return false
}
func (r *Router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	router := r.getRouter(path, request.Method)
	context := NewContext(response, request)

	if router == nil {
		router = &EmptyRouter
		context.NotFound()
		return
	}
	if err := request.ParseForm(); err != nil {
		logger.Error(fmt.Sprintf("parse form error: %s", err.Error()))
		context.BadRequest()
		return
	}
	if router.RequestType != nil {
		context.Input = reflect.New(router.RequestType).Interface()
		if isJsonType(request) {
			if err := json.NewDecoder(request.Body).Decode(context.Input); err != nil {
				logger.Error(fmt.Sprintf("json decode request data error: %s", err.Error()))
				context.BadRequest()
				return
			}
		} else {
			if err := formDecoder.Decode(context.Input, request.Form); err != nil {
				logger.Error(fmt.Sprintf("form decode request data error: %s", err.Error()))
				context.BadRequest()
				return
			}
		}
	}
	logger.Info(fmt.Sprintf("%s %s ,request data: %+v", request.Method, request.URL.Path, context.Input))
	for _, h := range router.Handler {
		h(context)
	}
	context.WriteToResponse()
	contextPool.Put(context)
}
func isJsonType(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Content-Type"), "application/json")
}
func (r *Router) PrintTree(path string) {
	if r.IsLeaf {
		if len(r.method) == 0 {
			r.method = append(r.method, "ALL")
		}
		for i := range r.method {
			logger.Info(fmt.Sprintf("%s %s\t%d handlers", r.method[i], path, len(r.Handler)))
			return
		}
	}
	for k := range r.Children {
		r.Children[k].PrintTree(path + "/" + k)
	}
}
