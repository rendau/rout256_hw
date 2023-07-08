package grpcserver

import (
	"fmt"
	"net"
	"route256/libs/logger"

	"google.golang.org/grpc"
)

type Server struct {
	Server *grpc.Server
	eChan  chan error
}

func New(opts ...grpc.ServerOption) *Server {
	return &Server{
		Server: grpc.NewServer(opts...),
		eChan:  make(chan error, 1),
	}
}

func (s *Server) Start(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("fail to listen: %w", err)
	}

	logger.Infow(nil, "Start grpc server", "addr", lis.Addr().String())

	go func() {
		err = s.Server.Serve(lis)
		if err != nil {
			logger.Errorw(nil, err, "GRPC server closed")
			s.eChan <- err
		}
	}()

	return nil
}

func (s *Server) Wait() <-chan error {
	return s.eChan
}

func (s *Server) Shutdown() bool {
	defer close(s.eChan)

	s.Server.GracefulStop()

	return true
}
