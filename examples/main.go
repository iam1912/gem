package main

import (
	"net/http"

	"github.com/iam1912/gem"
)

type student struct {
	Name string
	Age  int8
}

func main() {
	r := gem.New()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	stu1 := &student{Name: "Geektutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}

	r.GET("/", func(c *gem.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/students", func(c *gem.Context) {
		c.HTML(http.StatusOK, "home.html", gem.H{
			"title":  "gem",
			"stuArr": [2]*student{stu1, stu2},
		})
	})
	r.GET("/panic", func(c *gem.Context) {
		names := []string{"1111"}
		c.String(200, names[5])
	})
	// v1 := r.Group("/v1")
	// {
	// 	v1.GET("/", func(c *gem.Context) {
	// 		c.String(200, "HELLO WORLD")
	// 	})
	// 	v1.GET("/hello/:name", func(c *gem.Context) {
	// 		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	// 	})
	// }
	// v2 := r.Group("/v2")
	// v2.Use(onlyForV2())
	// {
	// 	v2.POST("/search", func(c *gem.Context) {
	// 		name := c.PostForm("name")
	// 		c.JSON(200, gem.H{
	// 			"name": name,
	// 		})
	// 	})
	// 	v2.GET("/assets/*filepath", func(c *gem.Context) {
	// 		c.JSON(http.StatusOK, gem.H{
	// 			"filepath": c.Param("filepath"),
	// 		})
	// 	})
	// }
	r.Run(":8080")
}
