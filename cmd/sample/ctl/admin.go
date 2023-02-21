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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// These are defined at the top-level command
var (
	abc string
)

// adminCmd is the management command
var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Sample CTL admin command",
	Long:  `The Sample CTL admin command will perform administrative operations.`,
	// PersistentPreRunE:
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugln("admin command called. Add main logic here")
		log.Debugf("%v", serverAddr)
		log.Infof("Starting adminCmd...")
	},
}

// init
func init() {
	//adminCmd.PersistenFlags().BoolVarP
	rootCmd.AddCommand(adminCmd)
}
