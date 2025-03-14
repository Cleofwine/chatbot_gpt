package server

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func GetTlsOpt(cert, key string) grpc.ServerOption {
	creds, err := credentials.NewServerTLSFromFile(cert, key)
	if err != nil {
		panic(err)
	}
	return grpc.Creds(creds)
}

func GetMTLSOpt(clientCaCert, certFile, keyFile string) grpc.ServerOption {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		panic(err)
	}
	ca := x509.NewCertPool()
	bytes, err := os.ReadFile(clientCaCert)
	if err != nil {
		panic(err)
	}
	ok := ca.AppendCertsFromPEM(bytes)
	if !ok {
		panic("append cert failed")
	}
	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    ca,
	}
	return grpc.Creds(credentials.NewTLS(tlsConfig))
}
