package grpc_protocol_plugin

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"github.com/zhengheng7913/grpc-go-starter/pkg/server"
	"log"
	"math/big"
	"reflect"
)

func init() {
	server.Register("quic", NewQuicService)
}

type Server interface {
	Serve(conn quic.Connection) error

	Close() error
}

func NewQuicService(opts ...server.Option) server.Service {
	options := &server.Options{}
	for _, f := range opts {
		f(options)
	}
	return &quicService{
		options: options,
	}
}

func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic"},
	}
}

type quicService struct {
	listener quic.Listener
	options  *server.Options
	svr      Server
}

func (s *quicService) Register(factory interface{}, impl interface{}) {
	s.svr = reflect.ValueOf(factory).Call([]reflect.Value{})[0].Interface().(Server)
}

func (s *quicService) Serve() (err error) {
	if s.listener, err = quic.ListenAddr(
		fmt.Sprintf("%s:%d", s.options.Host, s.options.Port),
		generateTLSConfig(),
		nil,
	); err != nil {
		return err
	}
	go s.listen()
	return nil
}

func (s *quicService) Close() error {
	if err := s.svr.Close(); err != nil {
		return err
	}
	if err := s.listener.Close(); err != nil {
		return err
	}
	return nil
}

func (s *quicService) listen() {
	for {
		conn, err := s.listener.Accept(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}
		go func() {
			if err := s.svr.Serve(conn); err != nil {
				log.Println(err)
			}
		}()
	}
}
