package server

import (
	"encoding/json"
	"html/template"
	"net"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Headers struct {
	Header string
	Value  string
}

type UserInfo struct {
	Headers []Headers
	IP      string
}

func getIP(r *http.Request) string {
	// Проверяем заголовок X-Forwarded-For (может содержать список IP)
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// Берём первый IP из списка
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}

	// Проверяем заголовок X-Real-IP
	ip = r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	// Используем RemoteAddr (формат: IP:port)
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // fallback
	}
	return ip
}

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", s.HelloWorldHandler)

	r.Get("/health", s.healthHandler)

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	headers := make([]Headers, 0)

	for header, value := range r.Header {
		headers = append(headers, Headers{
			Header: header,
			Value:  strings.Join(value, ", "),
		})
	}

	ip := getIP(r)

	userInfo := UserInfo{
		Headers: headers,
		IP:      ip,
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, userInfo)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
