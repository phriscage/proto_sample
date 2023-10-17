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
    "os"

    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
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
        log.Debugf("Starting bookCmd...")
        cmd.Help()
        log.Debugf("Finished bookCmd.")
        os.Exit(0)
    },
}

// init
func init() {
    bookCmd.PersistentFlags().StringVar(&bookName, "bookName", "", "Book name")
    rootCmd.AddCommand(bookCmd)
}
