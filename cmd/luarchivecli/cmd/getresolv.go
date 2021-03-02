// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/luids-io/api/dnsutil"
	dnsfinder "github.com/luids-io/api/dnsutil/grpc/finder"
)

// getresolvCmd represents the getresolv command
var getresolvCmd = &cobra.Command{
	Use:   "getresolv <id>...",
	Short: "Get resolv info using id",
	Long:  `Get resolv info using id`,

	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("id param is required")
		}
		return nil
	},

	Run: func(cmd *cobra.Command, args []string) {
		cli := dnsfinder.NewClient(grpcClient)
		ctx, cancel := getContextWithTimeout(context.Background())
		defer cancel()

		//prepare args
		jsonFormat, _ := cmd.Flags().GetBool("json")

		//process
		hasErrs := false
		for _, id := range args {
			data, found, err := cli.GetResolv(ctx, id)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %v\n", id, err)
				hasErrs = true
				continue
			}
			if !found {
				fmt.Printf("id: %s [not found]\n", id)
				continue
			}
			if jsonFormat {
				jsons, _ := json.Marshal(data)
				fmt.Printf("%s\n", string(jsons))
			} else {
				printResolv(os.Stdout, data)
			}
		}
		if hasErrs {
			os.Exit(1)
		}
	},
}

func printResolv(w io.Writer, r dnsutil.ResolvData) {
	fmt.Fprintf(w, "id: %s\n", r.ID)
	fmt.Fprintf(w, "timestamp: %s\n", r.Timestamp.Format(time.RFC3339))
	fmt.Fprintf(w, "server: %v\n", r.Server)
	fmt.Fprintf(w, "client: %v\n", r.Client)
	fmt.Fprintf(w, "qid: %v\n", r.QID)
	fmt.Fprintf(w, "name: %s\n", r.Name)
	fmt.Fprintf(w, "isIPv6: %v\n", r.IsIPv6)
	fmt.Fprintf(w, "queryFlags.authenticatedData: %v\n", r.QueryFlags.AuthenticatedData)
	fmt.Fprintf(w, "queryFlags.checkingDisabled: %v\n", r.QueryFlags.CheckingDisabled)
	fmt.Fprintf(w, "queryFlags.do: %v\n", r.QueryFlags.Do)
	fmt.Fprintf(w, "duration: %v\n", r.Duration)
	fmt.Fprintf(w, "returnCode: %v\n", r.ReturnCode)
	fmt.Fprintf(w, "responseFlags.authenticatedData: %v\n", r.ResponseFlags.AuthenticatedData)
	fmt.Fprintf(w, "resolvedIPs: %v\n", r.ResolvedIPs)
	fmt.Fprintf(w, "resolvedCNAMEs: %v\n", r.ResolvedCNAMEs)
	fmt.Fprintf(w, "tld: %s\n", r.TLD)
	fmt.Fprintf(w, "tldPlusOne: %s\n", r.TLDPlusOne)
}

func init() {
	rootCmd.AddCommand(getresolvCmd)

	getresolvCmd.Flags().Bool("json", false, "Json format")
}
