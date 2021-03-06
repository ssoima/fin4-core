package routes

import (
	"fin4-core/server/assethandlers"
	"fin4-core/server/datatype"
	"fin4-core/server/routermiddleware"
	"fin4-core/server/timelinehandlers"
	"fin4-core/server/tokenhandlers"
	"fin4-core/server/bysonhandlers"
	"github.com/gin-gonic/gin"



	"fin4-core/server/commonhandlers"
	"fin4-core/server/userhandlers"
	api "gopkg.in/appleboy/gin-status-api.v1"
)

//SetupRouting sets up routes
func SetupRouting(sc datatype.ServiceContainer) *gin.Engine {

	r := gin.Default()

	r.Static("/static", "./public")
	r.Static("/img", "./public/img")

	// html
	web := r.Group("/")
	web.Use(routermiddleware.Session())
	web.Use(routermiddleware.SessionSetUser(sc.UserService))
	{
		commonhandlers.InjectHandlers(sc, web)
	}

	// website specific api
	wapi := r.Group("/wapi")
	wapi.Use(routermiddleware.Session())
	wapi.Use(routermiddleware.SessionSetUser(sc.UserService))
	wapi.Use(routermiddleware.CheckCsrfToken())
	{
		wapi.GET("/csrf", routermiddleware.SetCsrfToken())
		userhandlers.InjectHandlers(sc, wapi)
		tokenhandlers.InjectHandlers(sc, wapi)
		bysonhandlers.InjectHandlers(sc, wapi)
		// V0.2 for assets
		assethandlers.InjectHandlers(sc, wapi)
		timelinehandlers.InjectHandlers(sc, wapi)
	}

	// API
	v1 := r.Group("/api")
	v1.Use(routermiddleware.HeadersNoCache())
	v1.Use(routermiddleware.HeadersCors())
	{
		v1.GET("/status", routermiddleware.APIAuth(), api.StatusHandler)
	}

	return r
}
