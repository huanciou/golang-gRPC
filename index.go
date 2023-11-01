package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"BLOG":   "www.flysnow.com",
			"wechat": "flysnow_org",
		})
	})

	r.GET("/hi", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"BLOG":   "www.hi.com",
			"wechat": "hi_org",
		})
	})
	r.Run(":8080")
}
