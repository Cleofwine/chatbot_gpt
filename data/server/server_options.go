package server

import "google.golang.org/grpc"

func GetOptions() (opts []grpc.ServerOption) {
	opts = make([]grpc.ServerOption, 0)
	opts = append(opts, GetKeepaliveOpt()...)
	// 考虑到data服务的是集群内部，不需要如此严格的鉴权方式
	// opts = append(opts, GetTlsOpt("./server/cert_ex/server_cert.pem", "./server/cert_ex/server_key.pem"))
	// opts = append(opts, GetMTLSOpt("./server/cert_ex/client_ca_cert.pem", "./server/cert_ex/server_cert.pem", "./server/cert_ex/server_key.pem"))
	opts = append(opts, grpc.StreamInterceptor(StreamInterceptor))
	opts = append(opts, grpc.UnaryInterceptor(UnaryInterceptor))
	return opts
}
