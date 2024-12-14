package cmd

import (
    "fmt"
    "os"
    "time"
    "github.com/spf13/cobra"
    "github.com/ohelal/rdap/pkg/rdap"
)

var (
    timeout  time.Duration
    baseURL  string
    verbose  bool
    format   string
    client   *rdap.Client
    outputStyle string
)

var rootCmd = &cobra.Command{
    Use:     "rdap",
    Short:   "RDAP CLI Tool - Query domain, IP, and ASN information",
    Version: "1.0.0",
    Long: `A command-line interface for querying RDAP (Registration Data Access Protocol) information.
Supports queries for domains, IP addresses, and Autonomous System Numbers (ASN).`,
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        // Setup client with options
        opts := []rdap.Option{rdap.WithTimeout(timeout)}
        if baseURL != "" {
            opts = append(opts, rdap.WithBaseURL(baseURL))
        }
        client = rdap.NewClient(opts...)
    },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func init() {
    rootCmd.PersistentFlags().DurationVarP(&timeout, "timeout", "t", 10*time.Second, "Query timeout")
    rootCmd.PersistentFlags().StringVar(&baseURL, "base-url", "", "Base URL for RDAP server")
    rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
    rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "pretty", "Output format: pretty, json, or compact")
    rootCmd.PersistentFlags().StringVarP(&outputStyle, "style", "s", "default", "Output style: default, table, box")
}
