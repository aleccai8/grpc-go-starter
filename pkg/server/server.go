package server

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	implementMap = make(map[string]func(opts ...Option) Service)
)

// Register 非线程安全
func Register(name string, constructor func(opts ...Option) Service) {
	implementMap[name] = constructor
}

func Get(name string) func(opts ...Option) Service {
	return implementMap[name]
}

func NewServer() *Server {
	return &Server{
		services: make(map[string]Service),
	}
}

type Server struct {
	services map[string]Service
}

func (s *Server) AddService(name string, service Service) {
	s.services[name] = service
}

func (s *Server) Service(name string) Service {
	return s.services[name]
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
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	_ = <-ch
	var wg sync.WaitGroup
	for _, service := range s.services {
		wg.Add(1)
		go func(srv Service) {
			defer wg.Done()
			if err := srv.Close(); err != nil {
				log.Println("server: close:", err)
			}
		}(service)
	}

	wg.Wait()
	return nil
}

func (s *Server) Close() error {
	return nil
}
