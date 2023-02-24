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
	configName string
)

// configCmd is the management command
var configCmd = &cobra.Command{
	Use:   "configs",
	Short: "Sample CTL configs command",
	Long:  `The Sample CTL configs command will perform configs operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debugf("Starting configCmd...")
		cmd.Help()
		log.Debugf("Finished configCmd")
		os.Exit(0)
	},
}

// init
func init() {
	configCmd.PersistentFlags().StringVar(&configName, "configName", "", "Config name")
	adminCmd.AddCommand(configCmd)
}
