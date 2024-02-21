package api

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"kitab/api/handler"
	"kitab/service"
)

func New(services service.IServiceManager) *gin.Engine {
	h := handler.New(services)

	r := gin.New()

	r.Use(traceRequest)

	r.POST("/book", h.CreateBook)
	r.GET("/book/:id", h.GetBook)
	r.GET("/books", h.GetBookList)
	r.PUT("/book/:id", h.UpdateBook)
	r.DELETE("/book/:id", h.DeleteBook)
	r.PATCH("/book/:id", h.UpdateBook)

	r.Use(afterRequest)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func traceRequest(c *gin.Context) {
	startTime := time.Now()

	c.Next()

	endTime := time.Now()
	duration := endTime.Sub(startTime).Milliseconds()

	log.Printf("%s %s status: %d, time: %d milliseconds\n", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration)
}

func beforeRequest(c *gin.Context) {
	startTime := time.Now()

	c.Set("start_time", startTime)

	log.Println("start time:", startTime.Format("2006-01-02 15:04:05.0000"), "path:", c.Request.URL.Path)
}

func afterRequest(c *gin.Context) {
	startTime, exists := c.Get("start_time")
	if !exists {
		startTime = time.Now()
	}

	duration := time.Since(startTime.(time.Time)).Seconds()

	fmt.Println("end time:", time.Now().Format("2006-01-02 15:04:05.0000"), "duration:", duration, "method:", c.Request.Method)
	fmt.Println()
}
