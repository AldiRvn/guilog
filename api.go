package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func mustJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}

	return string(b)
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func main() {
	appPort := ":8082"

	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()
	r.Use(func(ctx *gin.Context) {
		start := time.Now()

		reqBody, _ := io.ReadAll(ctx.Request.Body)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		writer := &responseWriter{
			ResponseWriter: ctx.Writer,
			body:           bytes.NewBuffer(nil),
		}

		ctx.Writer = writer

		ctx.Next()

		fmt.Printf(
			`logdy {"method":%q,"path":%q,"request_headers":%s,"request_body":%q,"status":%d,"response_headers":%s,"response_body":%q,"latency_ms":%d}`+"\n",
			ctx.Request.Method,
			ctx.Request.URL.RequestURI(),
			mustJSON(ctx.Request.Header),
			string(reqBody),
			ctx.Writer.Status(),
			mustJSON(ctx.Writer.Header()),
			writer.body.String(),
			time.Since(start).Milliseconds(),
		)
	})

	r.POST("/", func(ctx *gin.Context) {
		body := gin.H{}
		fmt.Printf("ctx.BindJSON(body): %v\n", ctx.BindJSON(&body))
		ctx.JSON(http.StatusOK, body)
	})

	// Define a simple GET endpoint
	r.GET("/", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	log.Println("http://0.0.0.0" + appPort)
	r.Run(appPort)
}
