package main

import (
	"log"
	"net/http"
	"time"

	"github.com/iam1912/gem"
)

type student struct {
	Name string
	Age  int8
}

func main() {
	r := gem.Default()

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	stu1 := &student{Name: "Geektutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}

	r.GET("/assets/*filepath", func(c *gem.Context) {
		c.JSON(http.StatusOK, gem.H{
			"filepath": c.Param("filepath"),
		})
	})
	r.GET("/hello/:name", func(c *gem.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})
	r.POST("/search", func(c *gem.Context) {
		name := c.PostForm("name")
		c.JSON(200, gem.H{
			"name": name,
		})
	})

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

	v1 := r.Group("/v1")
	v1.Use(onlyForV1())
	v1.GET("/test1", func(c *gem.Context) {
		c.String(200, "test1")
	})
	r.Run(":8082")
}

func onlyForV1() gem.HandlerFunc {
	return func(c *gem.Context) {
		t := time.Now()
		log.Printf("[%d] %s in %v for group v2 HELLO WORLD", c.StatusCode, c.Request.RequestURI, time.Since(t))
	}
}
