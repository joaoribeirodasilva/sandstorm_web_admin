package main

import (
	"github.com/joaoribeirodasilva/sandstorm_web_admin/webadmin/services/admin_log"
	"github.com/joaoribeirodasilva/sandstorm_web_admin/webadmin/services/config"
	"github.com/joaoribeirodasilva/sandstorm_web_admin/webadmin/services/ssl"
)

func main() {

	log := admin_log.New()

	log.Write("Starting Web Admin", "main", admin_log.LOG_INFO)

	config := config.New(log)
	config.Read()

	log.Open()

	ssl := ssl.New(config, log)
	ssl.Load()

	// router := gin.Default()

	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")

	// v1 := router.Group("/api")
	// {
	// 	v1.GET("/login", func(c *gin.Context) {
	// 		c.JSON(http.StatusOK, gin.H{"path": "/api/login"})
	// 	})
	// 	v1.GET("/submit", func(c *gin.Context) {
	// 		c.JSON(http.StatusOK, gin.H{"path": "/api/submit"})
	// 	})
	// 	v1.GET("/read", func(c *gin.Context) {
	// 		c.JSON(http.StatusOK, gin.H{"path": "/api/read"})
	// 	})
	// }

	// router.GET("/", func(c *gin.Context) {
	// 	router.LoadHTMLGlob("templates/*")
	// 	c.HTML(http.StatusOK, "index.html", gin.H{
	// 		"title": "Index",
	// 		"name":  "Home",
	// 	})
	// 	log.Printf("Path: %s\n", c.Request.URL.Path)
	// })
	// router.GET("/users", func(c *gin.Context) {
	// 	router.LoadHTMLGlob("templates/*")
	// 	c.HTML(http.StatusOK, "index.html", gin.H{
	// 		"title": "Users page",
	// 		"name":  "Users",
	// 	})
	// 	log.Printf("Path: %s\n", c.Request.URL.Path)
	// })
	// router.Run(":8080")
}
