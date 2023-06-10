package httpserver

import (
	"context"
	"log"
	"net/http"
	"time"
)

const (
	ReadHeaderTimeout = 2 * time.Second
	ReadTimeout       = 10 * time.Second
)

type Server struct {
	addr   string
	server *http.Server
	eChan  chan error
}

func Start(port string, handler http.Handler) *Server {
	s := &Server{
		addr: port,
		server: &http.Server{
			Addr:              ":" + port,
			Handler:           handler,
			ReadHeaderTimeout: ReadHeaderTimeout,
			ReadTimeout:       ReadTimeout,
		},
		eChan: make(chan error, 1),
	}

	log.Println("Start rest-api:", s.server.Addr)

	go func() {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Println("Http server closed", err)
			s.eChan <- err
		}
	}()

	return s
}

func (s *Server) Wait() <-chan error {
	return s.eChan
}

func (s *Server) Shutdown(timeout time.Duration) bool {
	defer close(s.eChan)

	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	defer ctxCancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		log.Println("Fail to shutdown http-api", "addr", s.addr, err)
		return false
	}

	return true
}
