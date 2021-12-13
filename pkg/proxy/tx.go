package proxy

import (
	"github.com/hyperledger-labs/yui-relayer/core"
)

func UpdateUpstreamClient(upstream *core.ProvableChain) error {
	pr := upstream.ProverI.(*Prover)
	// update the lightDB corresponding to the proxy
	if _, _, _, err := pr.upstreamProxy.UpdateLightWithHeader(); err != nil {
		return err
	}
	return pr.proxySynchronizer.SyncClientState()
}
