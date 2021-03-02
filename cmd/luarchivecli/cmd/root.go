// Copyright 2021 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/luids-io/common/config"
	"github.com/luids-io/common/factory"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var cfgFile string
var cfgClient config.ClientCfg
var cfgTimeoutSecs int

var grpcClient *grpc.ClientConn

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "luarchivecli",
	Short: "Archive API client",
	Long: `Archive API client.
luarchivecli is a command line implementation of the luarchive client API.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initialize)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.luarchivecli.yaml)")
	rootCmd.PersistentFlags().StringVarP(&cfgClient.RemoteURI, "uri", "r", "tcp://127.0.0.1:5821", "URI to grpc service.")
	rootCmd.PersistentFlags().StringVar(&cfgClient.TLS.CertFile, "clientcert", cfgClient.TLS.CertFile, "Path to grpc client cert file.")
	rootCmd.PersistentFlags().StringVar(&cfgClient.TLS.KeyFile, "clientkey", cfgClient.TLS.KeyFile, "Path to grpc client key file.")
	rootCmd.PersistentFlags().StringVar(&cfgClient.TLS.ServerCert, "servercert", cfgClient.TLS.ServerCert, "Path to grpc server cert file.")
	rootCmd.PersistentFlags().StringVar(&cfgClient.TLS.ServerName, "servername", cfgClient.TLS.ServerName, "Server name of grpc service for TLS check.")
	rootCmd.PersistentFlags().StringVar(&cfgClient.TLS.CACert, "cacert", cfgClient.TLS.CACert, "Path to grpc CA cert file.")
	rootCmd.PersistentFlags().BoolVar(&cfgClient.TLS.UseSystemCAs, "systemca", cfgClient.TLS.UseSystemCAs, "Use system CA pool for grpc check.")
	rootCmd.PersistentFlags().IntVar(&cfgTimeoutSecs, "timeout", 5, "Timeout in seconds for requests")
}

func exitWithErrf(format string, opts ...interface{}) {
	fmt.Fprintf(os.Stderr, format, opts...)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func initialize() {
	initConfig()
	initClient()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".tlsnotarycli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".luarchivecli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initClient() {
	err := cfgClient.Validate()
	if err != nil {
		exitWithErrf("bad grpc client config: %v", err)
	}
	con, err := factory.ClientConn(&cfgClient)
	if err != nil {
		exitWithErrf("creating grpc client connection: %v\n", err)
	}
	grpcClient = con
}

func getContextWithTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if cfgTimeoutSecs > 0 {
		return context.WithTimeout(ctx, time.Duration(cfgTimeoutSecs)*time.Second)
	}
	return context.WithCancel(ctx)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
