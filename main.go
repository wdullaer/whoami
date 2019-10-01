//go:generate protoc -I ./ --go_out=plugins=grpc:./pb/ ./whoami.proto

package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/prometheus/common/log"
	"github.com/wdullaer/whoami/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

var (
	port       = flag.String("port", os.Getenv("PORT"), "The port number to bind on (default: 50501)")
	portNumber = 50501
	tls        = flag.Bool("tls", false, "Set if the server should use TLS rather than plaintext")
	certFile   = flag.String("cert-file", "", "The path to the PEM encoded server certificate")
	keyFile    = flag.String("key-file", "", "The path to the PEM encoded key for the server certificate")
)

func main() {
	log.Info("WhomaiService starting")
	flag.Parse()
	if *port != "" {
		if portInt, err := strconv.Atoi(*port); err != nil {
			log.Fatalf("Failed to parse port number: %s", err)
		} else {
			portNumber = portInt
		}
	}

	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			log.Fatalf("TLS is requested but no certificate location was given")
		}
		if *keyFile == "" {
			log.Fatalf("TLS is requested but no key location for the certificate was given")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}

	whoamiServer := new(server)
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterWhoamiServiceServer(grpcServer, whoamiServer)
	reflection.Register(grpcServer)

	// Ensure a graceful shutdown on SIGTERM and SIGINT
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		event := <-signalChan
		log.Warnf("Received event to stop: %s", event.String())
		grpcServer.GracefulStop()
	}()

	// Start the gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", portNumber))
	if err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
	log.Infof("Starting gRPCServer on port: %d", portNumber)
	if err := grpcServer.Serve(lis); err != nil {
		log.Errorf("Dirty Shutdown: %s", err)
	}
	defer log.Warn("Shutdown completed")
}
