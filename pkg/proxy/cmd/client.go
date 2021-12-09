package cmd

import (
	"fmt"

	"github.com/datachainlab/ibc-proxy-relay/pkg/proxy"
	"github.com/hyperledger-labs/yui-relayer/config"
	"github.com/spf13/cobra"
)

func clientCmd(ctx *config.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-client [path] [upstream]",
		Short: "Update upstream client",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, src, dst, err := ctx.Config.ChainsFromPath(args[0])
			if err != nil {
				return err
			}
			switch args[1] {
			case src:
				return proxy.UpdateUpstreamClient(c[src])
			case dst:
				return proxy.UpdateUpstreamClient(c[dst])
			default:
				return fmt.Errorf("unknown chain %v", args[1])
			}
		},
	}
	return cmd
}
