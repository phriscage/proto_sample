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
    "os"

    "github.com/mitchellh/go-homedir"
    log "github.com/sirupsen/logrus"
    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var (
    VERSION      string
    cfgFile      string
    logSeverity  string
    tls          bool
    serverCAFile string
    serverAddr   string
)

// Bootstrap the Config object from env and defaults
//var cfg = &pb.Config{}

// rootCmd is the base command function when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "samplectl",
    Short: "Sample control",
    Long: `The Sample control command is an example for how to interface with protobuf 
messages and services via a CLI client.`,
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        if level, ok := log.ParseLevel(logSeverity); ok == nil {
            log.SetLevel(level)
        }
        requiredFlags := [1]string{"serverAddr"}
        for _, flagName := range requiredFlags {
            flagValue, _ := cmd.Flags().GetString(flagName)
            if flagValue == "" {
                log.Fatalf("'%s', requires a valid value and cannot be blank", flagName)
            }
        }
        return nil
    },
}

// Execute adds all the subcommands to the root command and sets the flags.
// It is called by main.main() and only needs to be set once to the rootCmd
func Execute(version string) {
    VERSION = version
    if err := rootCmd.Execute(); err != nil {
        log.Warn(err)
        os.Exit(1)
    }
}

func init() {
    // Log as JSON instead of the default ASCII formatter
    // log.SetFormatter(&log.JSONFormatter{})

    // Output to STDOUT instead of default STDERR
    log.SetOutput(os.Stdout)

    // Log severity set to Debug by default and will get changed
    log.SetLevel(log.DebugLevel)

    // Set the timestamp format
    log.SetFormatter(&log.TextFormatter{TimestampFormat: "2023-02-08T01:02:03.000000Z", FullTimestamp: true})

    cobra.OnInitialize(initConfig)

    // Define flags and configuration settings
    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.client.yaml)")

    // Verbose logging
    rootCmd.PersistentFlags().StringVarP(&logSeverity, "logSeverity", "l", getEnvOrString("LOG_SEVERITY", "INFO"), "Set the log severity. Defaults to env $LOG_SEVERITY")

    // Additional server parameters
    // Server Address
    rootCmd.PersistentFlags().StringVarP(&serverAddr, "serverAddr", "a", getEnvOrString("SERVER_ADDR", "127.0.0.1:10000"), "Set the server address to connect with. Defaults to env $SERVER_ADDR")

}

// initConfig reads in config file and environment variables if set
func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        home, err := homedir.Dir()
        if err != nil {
            log.Warn(err)
            os.Exit(1)
        }
        // Search in home directory
        viper.AddConfigPath(home)
        viper.SetConfigName(".client")
    }

    viper.AutomaticEnv() // read in envionrment variables that match

    if err := viper.ReadInConfig(); err == nil {
        fmt.Println("Using config file:", viper.ConfigFileUsed())
    }
}
