package router

import (
	"nanoCache/go-iml/cache"
	"nanoCache/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(c *cache.Cache) *gin.Engine {
	r := gin.Default()
	// r.Use()
	cache_group := r.Group("cache")
	{
		cache_group.POST("/set", controllers.CreateCacheHandler(c))
		cache_group.GET("/:key", controllers.GetCacheHandler(c))
		cache_group.DELETE("/:key", controllers.DeleteCacheHandler(c))
	}

	return r
}
