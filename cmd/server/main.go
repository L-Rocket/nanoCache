package main

import (
	"nanoCache/go-iml/cache"
	"nanoCache/internal/router"
)

func main() {
	c := cache.NewCache(256)
	r := router.SetupRouter(c)
	r.Run(":8080")
}
