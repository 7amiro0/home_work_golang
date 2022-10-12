package internalhttp

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

func createLog(ip, method, path, version, answer, browser string, latency time.Duration, date time.Time) string {
	return fmt.Sprintf("%s [%s] %s %s %s %s %s %s",
		ip, date.Format("02-01-2006 15:04:05"), method,
		path, version, answer, latency, browser)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		date := time.Now()
		method := r.Method
		path := r.URL.Path
		version := r.Proto
		answer := "200"
		browser := strings.Split(r.UserAgent(), " ")[0]

		latency := time.Since(start)
		l.Info(createLog(ip, method, path, version, answer, browser, latency, date))
	})
}
