package controllers

import (
	"nanoCache/go-iml/cache"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Request struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
	TTL   int    `json:"ttl" binding:"required"`
}

// set cache value with ttl, using post method
func CreateCacheHandler(cache *cache.Cache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Request
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		cache.Set(req.Key, req.Value, time.Duration(req.TTL)*time.Second)
		ctx.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}

// get cache value by key, use get method
func GetCacheHandler(cache *cache.Cache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.Param("key")
		value, ok := cache.Get(key)
		if !ok {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "key not found"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"key": key, "value": value})
	}
}

// delete cache value by key, use delete method
func DeleteCacheHandler(cache *cache.Cache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.Param("key")
		cache.Delete(key)
		ctx.JSON(http.StatusOK, gin.H{"status": "success"})
	}
}
