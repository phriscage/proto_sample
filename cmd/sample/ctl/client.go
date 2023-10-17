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
package ctl

import (
    "fmt"

    log "github.com/sirupsen/logrus"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/testdata"
    //pb "github.com/phriscage/proto_sample/gen/go/sample/v1alpha"
)

// newGRPCClientConn creates a new grpc.ClientConn to communicate with the GRPC server
func newGRPCClientConn(tls bool, srvCAFile, srvAddr string) (*grpc.ClientConn, error) {
    var opts []grpc.DialOption
    if tls {
        if srvCAFile == "" {
            srvCAFile = testdata.Path("ca.pem")
        }
        creds, err := credentials.NewClientTLSFromFile(srvCAFile, "")
        if err != nil {
            return nil, fmt.Errorf("FAiled to create TLS credentials &v", err)
        }
        opts = append(opts, grpc.WithTransportCredentials(creds))
    } else {
        opts = append(opts, grpc.WithInsecure())
    }

    log.Infof("Starting new gRPC client for server '%s'", srvAddr)
    conn, err := grpc.Dial(srvAddr, opts...)
    if err != nil {
        return nil, fmt.Errorf("failed to dial: %v", err)
    }
    return conn, nil
}
