package main

import (
	"log"

	proxy "github.com/datachainlab/ibc-proxy-prover/pkg/proxy/module"
	proxytm "github.com/datachainlab/ibc-proxy-prover/pkg/proxy/tendermint"
	tendermint "github.com/hyperledger-labs/yui-relayer/chains/tendermint/module"
	"github.com/hyperledger-labs/yui-relayer/cmd"
)

func main() {
	if err := cmd.Execute(
		tendermint.Module{},
		proxy.Module{},
		proxytm.Module{},
	); err != nil {
		log.Fatal(err)
	}
}
