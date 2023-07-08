package httpserver

import (
	"context"
	"net/http"
	"route256/libs/logger"
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

	logger.Infow(nil, "Start http server", "addr", s.server.Addr)

	go func() {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Errorw(nil, err, "Http server closed")
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
		logger.Errorw(nil, err, "Fail to shutdown http server", "addr", s.addr)
		return false
	}

	return true
}
