package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LoggingResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *LoggingResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Log request body for POST/PUT/PATCH
		var requestBody []byte
		if c.Request.Body != nil && (method == "POST" || method == "PUT" || method == "PATCH") {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Wrap response writer to capture response body
		w := &LoggingResponseWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get status code
		status := c.Writer.Status()

		// Log basic request info
		log.Printf("[%s] %s %s %d %v", 
			method, 
			path, 
			c.Request.Proto, 
			status, 
			duration,
		)

		// Log request body if present (for debugging)
		if len(requestBody) > 0 {
			var prettyJSON bytes.Buffer
			if json.Indent(&prettyJSON, requestBody, "", "  ") == nil {
				log.Printf("Request Body:\n%s", prettyJSON.String())
			}
		}

		// Log response body for errors
		if status >= 400 {
			log.Printf("Response Body (%d):\n%s", status, w.body.String())
		}

		// Log any errors
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				log.Printf("Error: %v", e.Error())
			}
		}
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}