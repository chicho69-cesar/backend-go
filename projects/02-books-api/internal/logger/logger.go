package logger

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Logger struct {
	file *os.File
}

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func NewLogger(logFilePath string) (*Logger, error) {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("Error al abrir archivo de log: %v", err)
	}

	logger := &Logger{file: file}

	logger.writeLog("=================================================")
	logger.writeLog(fmt.Sprintf("Logger iniciado - %s", time.Now().Format("2006-01-02 15:04:05")))
	logger.writeLog("=================================================")

	return logger, nil
}

func (l *Logger) writeLog(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] %s\n", timestamp, message)

	if l.file != nil {
		l.file.WriteString(logMessage)
	}

	log.Print(message)
}

func (l *Logger) LogRequest(method, path, remoteAddr, userAgent string, statusCode, responseSize int, duration time.Duration) {
	ip := remoteAddr
	if idx := strings.LastIndex(remoteAddr, ":"); idx != -1 {
		ip = remoteAddr[:idx]
	}

	logEntry := fmt.Sprintf(
		"%-7s | %-50s | IP: %-15s | Status: %3d | Size: %6d bytes | Time: %8s | UA: %s",
		method,
		path,
		ip,
		statusCode,
		responseSize,
		duration.String(),
		userAgent,
	)

	l.writeLog(logEntry)
}

func (l *Logger) Close() error {
	if l.file != nil {
		l.writeLog("=================================================")
		l.writeLog(fmt.Sprintf("Logger detenido - %s", time.Now().Format("2006-01-02 15:04:05")))
		l.writeLog("=================================================")

		return l.file.Close()
	}

	return nil
}

func (l *Logger) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := NewResponseWriter(w)

		next(rw, r)

		duration := time.Since(start)

		l.LogRequest(
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			r.UserAgent(),
			rw.statusCode,
			rw.size,
			duration,
		)
	}
}

func (l *Logger) MiddlewareHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := NewResponseWriter(w)

		next.ServeHTTP(rw, r)

		duration := time.Since(start)

		l.LogRequest(
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			r.UserAgent(),
			rw.statusCode,
			rw.size,
			duration,
		)
	})
}
