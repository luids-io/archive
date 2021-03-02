// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/luids-io/api/dnsutil"
	dnsfinder "github.com/luids-io/api/dnsutil/grpc/finder"
)

// getresolvCmd represents the getresolv command
var listresolvsCmd = &cobra.Command{
	Use:   "listresolvs",
	Short: "List resolvs",
	Long:  `List resolvs`,

	Run: func(cmd *cobra.Command, args []string) {
		cli := dnsfinder.NewClient(grpcClient)
		ctx, cancel := getContextWithTimeout(context.Background())
		defer cancel()

		//prepare args and filter
		rev, _ := cmd.Flags().GetBool("reverse")
		maxreq, _ := cmd.Flags().GetInt("maxreq")
		limit, _ := cmd.Flags().GetInt("limit")
		jsonFormat, _ := cmd.Flags().GetBool("json")
		f, err := getFilterFromFlags(cmd.Flags())
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		// do list
		var data []dnsutil.ResolvData
		next := ""
		count := 0
	LISTLOOP:
		for {
			data, next, err = cli.ListResolvs(ctx, []dnsutil.ResolvsFilter{f}, rev, maxreq, next)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
			for _, r := range data {
				if jsonFormat {
					jsons, _ := json.Marshal(r)
					fmt.Printf("%s\n", string(jsons))
				} else {
					fmt.Printf("%s,%s,%v,%s,%v,%v\n", r.ID, r.Timestamp.Format(time.RFC3339), r.Client, r.Name, r.ReturnCode, r.ResolvedIPs)
				}
				count++
				if limit > 0 && count >= limit {
					break LISTLOOP
				}
			}
			if next == "" {
				break
			}
		}
	},
}

func getFilterFromFlags(flags *pflag.FlagSet) (dnsutil.ResolvsFilter, error) {
	var err error
	var f dnsutil.ResolvsFilter
	cip, _ := flags.GetString("clientip")
	if cip != "" {
		f.Client = net.ParseIP(cip)
		if f.Client == nil {
			return f, errors.New("invalid 'clientip' format")
		}
	}
	sip, _ := flags.GetString("serverip")
	if sip != "" {
		f.Server = net.ParseIP(sip)
		if f.Server == nil {
			return f, errors.New("invalid 'serverip' format")
		}
	}
	f.Name, _ = flags.GetString("name")
	f.QID, _ = flags.GetInt("qid")
	f.ReturnCode, _ = flags.GetInt("returncode")
	rip, _ := flags.GetString("resolvedip")
	if rip != "" {
		f.ResolvedIP = net.ParseIP(rip)
		if f.ResolvedIP == nil {
			return f, errors.New("invalid 'resolvedip' format")
		}
	}
	f.ResolvedCNAME, _ = flags.GetString("resolvedcname")
	f.TLD, _ = flags.GetString("tld")
	f.TLDPlusOne, _ = flags.GetString("tldplusone")
	since, _ := flags.GetString("since")
	if since != "" {
		f.Since, err = time.Parse(time.RFC3339, since)
		if err != nil {
			return f, fmt.Errorf("invalid 'since' format: %v", err)
		}
	}
	to, _ := flags.GetString("to")
	if to != "" {
		f.To, err = time.Parse(time.RFC3339, to)
		if err != nil {
			return f, fmt.Errorf("invalid 'to' format: %v", err)
		}
	}
	return f, nil
}

func init() {
	rootCmd.AddCommand(listresolvsCmd)

	listresolvsCmd.Flags().Bool("reverse", false, "Reverse order")
	listresolvsCmd.Flags().Int("maxreq", 0, "Max items per fetch request")
	listresolvsCmd.Flags().Int("limit", 0, "Max items listed")
	listresolvsCmd.Flags().Bool("json", false, "Json format")
	//filter args
	listresolvsCmd.Flags().String("clientip", "", "Filter by client IP")
	listresolvsCmd.Flags().String("serverip", "", "Filter by server IP")
	listresolvsCmd.Flags().String("name", "", "Filter by name")
	listresolvsCmd.Flags().Int("qid", 0, "Filter by qid")
	listresolvsCmd.Flags().Int("returncode", 0, "Filter by return code")
	listresolvsCmd.Flags().String("resolvedip", "", "Filter by resolved IP")
	listresolvsCmd.Flags().String("resolvedcname", "", "Filter by resolved cname")
	listresolvsCmd.Flags().String("tld", "", "Filter by tld")
	listresolvsCmd.Flags().String("tldplusone", "", "Filter by tldplusone")
	listresolvsCmd.Flags().String("since", "", "Filter since timestamp (format '"+time.RFC3339+"')")
	listresolvsCmd.Flags().String("to", "", "Filter to timestamp (format '"+time.RFC3339+"')")
}
