/*
MIT License

# Copyright (c) 2023 phriscage

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "github.com/phriscage/proto_sample/gen/go/sample/v1alpha"
)

var (
	tls         = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	certFile    = flag.String("cert_file", "", "The TLS cert file")
	keyFile     = flag.String("key_file", "", "The TLS key file")
	port        = flag.Int("port", 10000, "The server port")
	host        = flag.String("host", "127.0.0.1", "The server host ip")
	logSeverity = flag.String("l", "INFO", "Set the log severity")
	environment = flag.String("e", "development", "Set the environment name")
)

// Sample Server object that includes the configurations
type sampleServer struct {
	// GRPC server
	pb.UnimplementedSampleServiceServer

	// Server Config
	ServerCfg *pb.Config

	// TODO Setup CSP configs
	//// GCP Clients
	// gcs *storage.Client

	//// AWS Clients
	// ebs *ebs.Client

	mu sync.Mutex // protects books

	// Collection of books
	books map[string][]*pb.Book
}

// Init
func init() {
	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)

	// Only log the debug severity or above
	log.SetLevel(log.DebugLevel)

	// Set the timestamp format in output
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2023-02-08T01:02:03.000000Z", FullTimestamp: true})
}

func defaultServerOpts() []grpc.ServerOption {
	return []grpc.ServerOption{}
}

// withDuration returns the duration of a grpc connection in nanoseconds
func withDuration(duration time.Duration) (key string, value interface{}) {
	return "grpc.time_ns", duration.Nanoseconds()
}

// Sample Server Get ServerCfg getter function
func (x *sampleServer) getServerCfg() *pb.Config {
	if x != nil {
		return x.ServerCfg
	}
	return nil
}

// Validate the Sample Server Config pb
func (x *sampleServer) validateServerCfg() error {
	// validate pb.Config
	serverCfg := x.getServerCfg()
	if serverCfg == nil {
		return fmt.Errorf("pb.Config cannot be nil")
	}
	/*
		// validate pb.Config.XyZ enum
		sev := serverCfg.GetXyZ()
		if _, ok := pb.Config_XyZ_value[sev.String()]; !ok {
			return fmt.Errorf("pb.Config.XyZ is not a valid value")
		}
	*/
	return nil
}

// Init a new Sample Server object and any downstream clients
func newSampleServer(serverCfg *pb.Config) *sampleServer {
	// Init the Sample Server
	s := &sampleServer{
		ServerCfg: serverCfg,
	}
	// Validate the Sample Server Config
	log.Debugf("Validating the Server Configs...")
	if err := s.validateServerCfg(); err != nil {
		log.Fatal(err)
	}
	log.Debugf("%+v", s)
	// TODO Setup the Sample Server client contexts
	//ctx := context.Background()
	return s
}

func main() {
	flag.Parse()
	opts := []grpc.ServerOption{
		// The following grpc.ServerOption adds an interceptor for all unary
		// RPCs. To configure an interceptor for streaming RPCs, see:
		// https://godoc.org/google.golang.org/grpc#StreamInterceptor
		//grpc.UnaryInterceptor(ensureValidToken),
		// Enable TLS for all incoming connections.
		//grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
	}
	if level, ok := log.ParseLevel(*logSeverity); ok == nil {
		log.Info(log.ParseLevel(*logSeverity))
		log.SetLevel(level)
	}
	host_port := fmt.Sprintf("%s:%d", *host, *port)
	lis, err := net.Listen("tcp", host_port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	/* TODO add TLS cert
	if *tls {
		if *certFile == "" {
			*certFile = data.Path("x509/server_cert.pem")
		}
		if *keyFile == "" {
			*keyFile = data.Path("x509/server_key.pem")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials: %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	*/
	logrusEntry := log.NewEntry(log.StandardLogger())
	logOpts := []grpc_logrus.Option{
		grpc_logrus.WithDurationField(withDuration),
	}
	opts = append(
		opts,
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logrusEntry, logOpts...)),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_logrus.StreamServerInterceptor(logrusEntry, logOpts...)),
	)

	log.Infof("Starting grpc server on '%s'", host_port)
	grpcServer := grpc.NewServer(append(defaultServerOpts(), opts...)...)

	cfg := &pb.Config{
		Environment: *environment,
	}
	//pb.RegisterSampleServiceServer(grpcServer)
	pb.RegisterSampleServiceServer(grpcServer, newSampleServer(cfg))
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}

	log.Info("Stoping grpc serverr...")
}
