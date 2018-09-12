package igin

import (
	"net/http"

	"os"

	"github.com/gin-gonic/gin"
	"study/stu/go_utils/reqid.v1"
)

type Server struct {
	*gin.Engine
	s *http.Server
}

func NewServer(addr ...string) *Server {
	r := gin.New()
	r.Use(gin.Recovery())
	handler := reqid.RequestId(r)

	address := resolveAddress(addr)
	return &Server{
		Engine: r,
		s: &http.Server{
			Addr:    address,
			Handler: handler,
		},
	}
}

func (s *Server) Run() error {
	return s.s.ListenAndServe()
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); len(port) > 0 {
			return ":" + port
		}
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too much parameters")
	}
}
