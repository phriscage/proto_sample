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

	pb "github.com/phriscage/proto_sample/gen/go/sample/v1alpha"
)

// These are defined at the top-level command
var (
	bookName string
)

// bookCmd is the management command
var bookCmd = &cobra.Command{
	Use:   "books",
	Short: "Sample CTL books command",
	Long:  `The Sample CTL books command will perform books operations.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		requiredFlags := [1]string{"bookName"}
		for _, flagName := range requiredFlags {
			flagValue, _ := cmd.Flags().GetString(flagName)
			if flagValue == "" {
				log.Fatalf("'%s', requires a valid value and cannot be blank", flagName)
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("Starting booksCmd...")

		conn, err := newGRPCClientConn(tls, serverCAFile, serverAddr)
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		c := pb.NewSampleServiceClient(conn)
		book, err := getBook(c, &pb.GetBookRequest{Name: bookName})
		if err != nil {
			log.Fatal(err)
		}
		log.Debugf("%#v", book)

		log.Infof("Finished getBooks")
	},
}

// getBook gets the Book from the server
func getBook(c pb.SampleServiceClient, req *pb.GetBookRequest) (*pb.Book, error) {
	log.Debugf("Getting Book for (%s)", req)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var header, trailer metadata.MD // var to store header and trailer
	resp, err := c.GetBook(ctx, req, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		log.Warnf("%#v.getBook(_) = _, %v: ", c, err)
		return nil, err
	}
	return resp, nil
}

// init
func init() {
	bookCmd.PersistentFlags().StringVar(&bookName, "bookName", "", "Book name")
	rootCmd.AddCommand(bookCmd)
}
