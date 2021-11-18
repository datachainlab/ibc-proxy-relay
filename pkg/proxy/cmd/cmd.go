package cmd

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/hyperledger-labs/yui-relayer/config"
	"github.com/spf13/cobra"
)

func ProxyCmd(m codec.Codec, ctx *config.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxy",
		Short: "Proxy commands",
	}

	cmd.AddCommand(
		clientCmd(ctx),
	)

	return cmd
}
