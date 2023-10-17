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
    "context"
    "flag"
    "fmt"
    "net"
    "os"
    "sync"
    "time"

    "github.com/davecgh/go-spew/spew"
    grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
    grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
    grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
    log "github.com/sirupsen/logrus"
    "golang.org/x/exp/maps"
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "google.golang.org/protobuf/types/known/emptypb"

    // TODO Using GROM v2 requires protoc-gen-gorm to update
    // dependencies [here](https://github.com/infobloxopen/protoc-gen-gorm/issues/243)
    // "gorm.io/driver/sqlite"
    // "gorm.io/gorm"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"

    pb "github.com/phriscage/proto_sample/gen/go/sample/v1alpha"
)

var (
    tls                   = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
    certFile              = flag.String("cert_file", getEnvOrString("CERT_FILE", ""), "The TLS cert file")
    keyFile               = flag.String("key_file", getEnvOrString("KEY_FILE", ""), "The TLS key file")
    port                  = flag.Int("port", 10000, "The server port")
    host                  = flag.String("host", getEnvOrString("HOST", "127.0.0.1"), "The server host ip")
    logSeverity           = flag.String("log_severity", getEnvOrString("LOG_SEVERITY", "INFO"), "Set the log severity")
    environment           = flag.String("environment", getEnvOrString("ENVIRONMENT", "development"), "Set the environment name")
    databaseProvider      = flag.String("database_provider", getEnvOrString("DATABASE_PROVIDER", "sqlite3"), "Set the Database provider type")
    databaseConnectionDsn = flag.String("database_connection_dsn", getEnvOrString("DATABASE_CONNECTION_DSN", "abc://123"), "Set the Database Connection DSN")
)

// Sample Server object that includes the configurations
type sampleServer struct {
    // GRPC server
    pb.UnimplementedSampleServiceServer

    // Server Config
    cfg *pb.Config

    // Database Client
    db *gorm.DB

    // TODO Setup CSP configs
    //// GCP Clients
    // gcs *storage.Client

    //// AWS Clients
    // ebs *ebs.Client

    mu sync.Mutex // protects books

    // Collection of books
    books map[string]*pb.Book
}

//
//    helper functions for the grpcServer
//

// defaultServer options
func defaultServerOpts() []grpc.ServerOption {
    return []grpc.ServerOption{}
}

// withDuration returns the duration of a grpc connection in nanoseconds
func withDuration(duration time.Duration) (key string, value interface{}) {
    return "grpc.time_ns", duration.Nanoseconds()
}

// Init a new Sample Server object and any downstream clients
func newSampleServer(cfg *pb.Config) *sampleServer {
    // Init the Sample Server
    s := &sampleServer{
        cfg:   cfg,
        books: make(map[string]*pb.Book),
    }
    // Validate the Sample Server Config
    log.Debugf("Validating the Server Configs...")
    if err := s.validateCfg(); err != nil {
        log.Fatal(err)
    }
    log.Debug(spew.Sprintf("%#v", s))
    // Setup database connection(s)
    //log.Debug(s.getCfg())
    dsn := s.getCfg().GetDatabase().GetConnection().GetDsn()
    provider := s.getCfg().GetDatabase().GetProvider()
    var db *gorm.DB
    var err error
    if provider.String() == "SQLITE" {
        // db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{}) // GORM v2
        db, err = gorm.Open("sqlite3", dsn)
        // TODO add PostgreSQL support
        //} else if provider.String() == "POSTGRESQL" {
    } else {
        log.Fatalf("pb.Config.Database.Provider not configured")
    }
    if err != nil {
        log.Fatal(err)
    }
    s.db = db
    db.LogMode(false)
    // TODO Setup the Sample Server client contexts
    //ctx := context.Background()
    return s
}

// Sample Server Get Cfg getter function
func (x *sampleServer) getCfg() *pb.Config {
    if x != nil {
        return x.cfg
    }
    return nil
}

// Validate the Sample Server Config pb
func (x *sampleServer) validateCfg() error {
    // validate pb.Config
    cfg := x.getCfg()
    if cfg == nil {
        return fmt.Errorf("pb.Config cannot be nil")
    }
    /*
        // validate pb.Config.XyZ enum
        sev := cfg.GetXyZ()
        if _, ok := pb.Config_XyZ_value[sev.String()]; !ok {
            return fmt.Errorf("pb.Config.XyZ is not a valid value")
        }
    */
    dbCfg := cfg.GetDatabase()
    if dbCfg == nil {
        return fmt.Errorf("pb.Config.Database cannot be nil")
    }
    if dsn := dbCfg.GetConnection().GetDsn(); dsn == "" {
        return fmt.Errorf("pb.Config.Database.Connection.Dsn cannot be nil")
    }

    return nil
}

//
// grpc server methods
//

// GetConfig method
func (s *sampleServer) GetConfig(ctx context.Context, _ *emptypb.Empty) (*pb.Config, error) {
    log.Infof("Starting GetConfig...")
    if s.getCfg() == nil {
        return &pb.Config{}, status.Error(codes.NotFound, fmt.Sprintf("Does not exist"))
    }
    return s.getCfg(), status.Error(codes.OK, codes.OK.String())
}

// GetBook method
func (s *sampleServer) GetBook(ctx context.Context, req *pb.GetBookRequest) (*pb.Book, error) {
    log.Infof("Starting GetBook...")
    if req == nil && req.GetName() == "" {
        return &pb.Book{}, status.Error(codes.InvalidArgument, fmt.Sprintf("Request is not valid"))
    }
    provider := s.getCfg().GetDatabase().GetProvider()
    pbBook := pb.Book{Name: req.GetName()}
    book, err := pbBook.ToORM(ctx)
    if err != nil {
        log.Errorf("%s: %s", provider, err)
        return &pb.Book{}, status.Error(codes.Internal, fmt.Sprint("Internal Server Error"))
    }
    if err := s.db.Where(&pbBook).First(&book).Error; err != nil {
        log.Warnf("%s: %s", provider, err)
        return &pb.Book{}, status.Error(codes.NotFound, fmt.Sprintf("Does not exist"))
    }
    pbBook, err = book.ToPB(ctx)
    if err != nil {
        log.Errorf("%s: %s", provider, err)
        return &pb.Book{}, status.Error(codes.Internal, fmt.Sprint("Internal Server Error"))
    }
    return &pbBook, status.Error(codes.OK, codes.OK.String())
}

// CreateBook method
func (s *sampleServer) CreateBook(ctx context.Context, req *pb.CreateBookRequest) (*pb.CreateBookResponse, error) {
    log.Infof("Starting CreateBook...")
    // Name, authors [Person], title
    if req == nil || req.GetBook() == nil {
        return &pb.CreateBookResponse{StatusMessage: "Request is not valid"}, status.Error(codes.InvalidArgument, fmt.Sprintf("Request is not valid"))
    }
    provider := s.getCfg().GetDatabase().GetProvider()
    pbBook := req.GetBook()
    book, err := pbBook.ToORM(ctx)
    book.Id = createUUIDv4()
    if err != nil {
        log.Errorf("%s: %s", provider, err)
        return &pb.CreateBookResponse{StatusMessage: codes.Internal.String()}, status.Error(codes.Internal, fmt.Sprint("Internal Server Error"))
    }
    if err := s.db.Where(&pbBook).First(&book).Error; err == nil {
        return &pb.CreateBookResponse{StatusMessage: codes.AlreadyExists.String()}, status.Error(codes.AlreadyExists, codes.AlreadyExists.String())
    }
    log.Debug(spew.Sprintf("%#v", book))
    if err := s.db.Create(&book).Error; err != nil {
        log.Errorf("%s: %s", provider, err)
        return &pb.CreateBookResponse{StatusMessage: codes.Internal.String()}, status.Error(codes.Internal, fmt.Sprint("Internal Server Error"))
    }
    return &pb.CreateBookResponse{StatusMessage: codes.OK.String()}, status.Error(codes.OK, codes.OK.String())
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

// main
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

    // Config init and validation
    //var dbProviderVal pb.Database_Provider
    dbProviderVal, ok := pb.Database_Provider_value[*databaseProvider]
    if !ok {
        providers := remove(maps.Values(pb.Database_Provider_name), "UNKNOWN")
        log.Fatalf("pb.Config.Database.Provider is not a valid name: '%s'", providers)
    }
    cfg := &pb.Config{
        Environment: *environment,
        Database: &pb.Database{
            Provider: pb.Database_Provider(dbProviderVal),
            Connection: &pb.Database_Connection{
                Dsn: *databaseConnectionDsn,
            },
        },
    }
    //pb.RegisterSampleServiceServer(grpcServer)
    pb.RegisterSampleServiceServer(grpcServer, newSampleServer(cfg))
    if err = grpcServer.Serve(lis); err != nil {
        log.Fatal(err)
    }

    log.Info("Stoping grpc serverr...")
}
