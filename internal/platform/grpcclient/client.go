package grpcclient

import (
	"fmt"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewInsecure(addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return newInsecure(strings.TrimSpace(addr), opts...)
}

func NewRequiredInsecure(name string, addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "gRPC target"
	}
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return nil, fmt.Errorf("%s address is required", name)
	}
	conn, err := NewInsecure(addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("connect to %s: %w", name, err)
	}
	return conn, nil
}

func NewInsecurePassthrough(addr string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return newInsecure(TargetPassthrough(addr), opts...)
}

func SelfTarget(listenAddr string) string {
	addr := strings.TrimSpace(listenAddr)
	if strings.HasPrefix(addr, ":") {
		return "127.0.0.1" + addr
	}
	host, port, err := net.SplitHostPort(addr)
	if err != nil || port == "" {
		return addr
	}
	if host == "" || host == "0.0.0.0" || host == "::" {
		host = "127.0.0.1"
	}
	return net.JoinHostPort(host, port)
}

func TargetPassthrough(addr string) string {
	addr = strings.TrimSpace(addr)
	if addr == "" || strings.Contains(addr, "://") || strings.HasPrefix(addr, "passthrough:") {
		return addr
	}
	return "passthrough:///" + addr
}

func newInsecure(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	options := append([]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}, opts...)
	return grpc.NewClient(target, options...)
}
