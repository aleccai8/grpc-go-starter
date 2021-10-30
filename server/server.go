package server

import (
	"context"
	"os"
	"sync"
	"syscall"
	"time"
)

const MaxCloseWaitTime = 10 * time.Second

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

func (s *Server) Register(desc interface{}, serviceImpl interface{}) {
	for _, srv := range s.services {
		srv.Register(desc, serviceImpl)
	}
}

func (s *Server) Serve() error {
	if len(s.services) == 0 {
		panic("service empty")
	}

	ch := make(chan os.Signal)

	for name, service := range s.services {
		go func(n string, srv Service) {
			if e := srv.Serve(); e != nil {
				ch <- syscall.SIGTERM
				panic(e)
			}
		}(name, service)
	}
	_ = <-ch

	ctx, cancel := context.WithTimeout(context.Background(), MaxCloseWaitTime)
	defer cancel()
	var wg sync.WaitGroup
	for _, service := range s.services {

		wg.Add(1)
		go func(srv Service) {
			defer wg.Done()

			c := make(chan struct{}, 1)
			go srv.Close(c)
			select {
			case <-c:
			case <-ctx.Done():
			}
		}(service)
	}

	wg.Wait()
	return nil
}

func (s *Server) Close(chan struct{}) error {
	return nil
}
