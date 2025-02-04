package router

import (
	"github.com/gin-gonic/gin"
)

type HttpRequestType int

const (
	Get HttpRequestType = iota
	Post
	Put
	Patch
	Delete
)

var routes = map[string]map[HttpRequestType]func(*gin.Context){
	"/checkAuth": {},
}

func InitHttpRoutes(router gin.IRoutes) {
	for route, requestTypePairs := range routes {
		for requestType, function := range requestTypePairs {
			switch requestType {
			case Get:
				router.GET(route, function)

			case Post:
				router.POST(route, function)

			case Put:
				router.PUT(route, function)

			case Patch:
				router.PATCH(route, function)

			case Delete:
				router.DELETE(route, function)
			}
		}
	}
}
