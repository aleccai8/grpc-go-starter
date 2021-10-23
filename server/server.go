package server

import "fmt"

func NewServer() *Server {
	return &Server{
		services: make(map[string]Service),
	}
}

type Server struct {
	services map[string]Service
}

func (s *Server) AddService(serviceName string, service Service) {
	s.services[serviceName] = service
}

func (s *Server) Service(serviceName string) Service {
	return s.services[serviceName]
}

func (s *Server) Register(serviceDesc *ServiceDesc, serviceImpl interface{}) error {
	return fmt.Errorf("can not register server as service")
}

func (s *Server) Serve() error {
	return nil
}

func (s *Server) Close(chan struct{}) error {
	return nil
}
