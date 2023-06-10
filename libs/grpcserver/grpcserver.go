package grpcserver

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	Server *grpc.Server
	eChan  chan error
}

func New() *Server {
	return &Server{
		Server: grpc.NewServer(),
		eChan:  make(chan error, 1),
	}
}

func (s *Server) Start(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("fail to listen: %w", err)
	}

	log.Println("Start grpc server:", lis.Addr().String())

	go func() {
		err = s.Server.Serve(lis)
		if err != nil {
			log.Println("GRPC server closed", err)
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
