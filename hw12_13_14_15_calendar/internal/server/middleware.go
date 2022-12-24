package server

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

func createLog(ip, method, path, version, answer string, browser []string, latency time.Duration, date time.Time) string {
	return fmt.Sprintf("%s [%s] %s %s %s %s %s %s",
		ip, date.Format("02-01-2006 15:04:05"), method,
		path, version, answer, latency, browser)
}

func (middle) LoggingMiddleware(next http.Handler, logger Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			logger.Fatal(err)
		}

		date := time.Now()
		method := r.Method
		path := r.URL.Path
		version := r.Proto
		answer := "200"
		browser := strings.Split(r.UserAgent(), " ")
		latency := time.Since(start)
		logger.Info(createLog(ip, method, path, version, answer, browser, latency, date))
	})
}
