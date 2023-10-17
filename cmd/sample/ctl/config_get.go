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
    "context"
    "time"

    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
    "google.golang.org/protobuf/types/known/emptypb"

    pb "github.com/phriscage/proto_sample/gen/go/sample/v1alpha"
)

// configGetCmd is the management command
var configGetCmd = &cobra.Command{
    Use:   "get",
    Short: "Sample CTL config get command",
    Long:  `The Sample CTL configs command will perform config get operation.`,
    Run: func(cmd *cobra.Command, args []string) {
        log.Infof("Starting configGetCmd...")

        conn, err := newGRPCClientConn(tls, serverCAFile, serverAddr)
        if err != nil {
            log.Fatal(err)
        }
        defer conn.Close()
        c := pb.NewSampleServiceClient(conn)
        config, err := getConfig(c, &emptypb.Empty{})
        if err != nil {
            log.Fatal(err)
        }
        log.Debugf("%#v", config)

        log.Infof("Finished configGetCmd.")
    },
}

// getConfig gets the Config from the server
func getConfig(c pb.SampleServiceClient, req *emptypb.Empty) (*pb.Config, error) {
    log.Debugf("Getting Config for (%s)", req)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    var header, trailer metadata.MD // var to store header and trailer
    resp, err := c.GetConfig(ctx, req, grpc.Header(&header), grpc.Trailer(&trailer))
    if err != nil {
        log.Warnf("%#v.getConfig(_) = _, %v: ", c, err)
        return nil, err
    }
    return resp, nil
}

// init
func init() {
    //configsCmd.PersistenFlags().BoolVarP
    configCmd.AddCommand(configGetCmd)
}
