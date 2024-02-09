package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type HandlerInfo struct {
	HandlerFunc http.HandlerFunc
	Method      string
}

type WebServer struct {
	Router   chi.Router
	Handlers map[string]HandlerInfo
	Port     string
}

func NewWebServer(port string) *WebServer {
	return &WebServer{
		Router:   chi.NewRouter(),
		Handlers: make(map[string]HandlerInfo),
		Port:     port,
	}
}

func (s *WebServer) AddHandler(path, method string, handler http.HandlerFunc) {
	s.Handlers[path] = HandlerInfo{
		Method:      strings.ToUpper(method),
		HandlerFunc: handler,
	}
}

func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)

	for path, handler := range s.Handlers {
		s.Router.Handle(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != handler.Method {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
				return
			}
			handler.HandlerFunc.ServeHTTP(w, r)
		}))
	}

	fmt.Println("App is running on http://localhost:" + s.Port)

	http.ListenAndServe(":"+s.Port, s.Router)
}
